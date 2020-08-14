package basecell

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/davyxu/golog"

	"github.com/adamluo159/cellnetEx"
)

var log *golog.Logger = nil //= golog.New("websocket_bceller")
//DefaultCell 默认服务
var DefaultCell *BaseCell = nil

//IModule 模块接口
type iModule interface {
	Init()
	Name() string
	OnDestory()
}

//iUserData 用户数据接口
type iUserData interface {
	QID() int
	UID() string
}

//BaseCell 基础服务
type BaseCell struct {
	MsgQueueLen int

	//bcellName    string //服务名字
	modules    []iModule
	msgHandler map[reflect.Type]func(ev cellnetEx.Event)
	queue      cellnetEx.EventQueue
	queues     []cellnetEx.EventQueue

	peer cellnetEx.GenericPeer
}

//SetLog 设置日志
func SetLog(l *golog.Logger) {
	log = l
}

//New 创建新服务
func New(msgQueLen int) *BaseCell {
	if msgQueLen < 0 {
		panic("msgQueLen < 0")
	}

	if msgQueLen%2 == 0 && msgQueLen > 0 {
		panic("need msgQueLen % 2 != 0")
	}

	bcell := &BaseCell{
		MsgQueueLen: msgQueLen,
		queue:       cellnetEx.NewEventQueue(),
		queues:      make([]cellnetEx.EventQueue, 0),
		msgHandler:  make(map[reflect.Type]func(ev cellnetEx.Event)),
	}

	bcell.queue.EnableCapturePanic(true)
	for i := 0; i < msgQueLen; i++ {
		q := cellnetEx.NewEventQueue()
		q.EnableCapturePanic(true)
		bcell.queues = append(bcell.queues, q)
	}

	if DefaultCell == nil {
		DefaultCell = bcell
	}
	return bcell
}

func (bcell *BaseCell) msgQueue() func(ev cellnetEx.Event) {
	return func(ev cellnetEx.Event) {
		if bcell.MsgQueueLen > 0 {
			queueID := 0
			udata := ev.Session().GetUserData()
			if udata == nil {
				queueID = rand.Intn(bcell.MsgQueueLen)
			} else {
				queueID = udata.(iUserData).QID()
			}
			bcell.queues[queueID].Post(func() {
				f, ok := bcell.msgHandler[reflect.TypeOf(ev.Message())]
				if ok {
					f(ev)
				} else {
					log.Errorln("onMessage not found message handler ", ev.Message())
				}
			})
			return
		}

		f, ok := bcell.msgHandler[reflect.TypeOf(ev.Message())]
		if ok {
			f(ev)
		} else {
			log.Errorln("onMessage not found message handler ", ev.Message())
		}
	}
}

//Start 服务开始
func (bcell *BaseCell) Start(mods ...iModule) {
	tmpNames := []string{}
	for _, m := range mods {
		for _, name := range tmpNames {
			if name == m.Name() {
				panic(fmt.Sprintf("repeat module name:%s", m.Name()))
			}
		}
		m.Init()
		tmpNames = append(tmpNames, m.Name())
	}
	bcell.modules = mods
	// 开始侦听
	bcell.peer.Start()

	// 事件队列开始循环
	bcell.queue.StartLoop()

	for _, v := range bcell.queues {
		v.StartLoop()
	}
}

//Stop 服务停止
func (bcell *BaseCell) Stop() {
	bcell.peer.Stop()
	bcell.queue.StopLoop()
	bcell.queue.Wait()

	for _, v := range bcell.queues {
		v.StopLoop()
		v.Wait()
	}

	for _, m := range bcell.modules {
		m.OnDestory()
	}
}

//RegisterMessage 注册默认消息响应
func RegisterMessage(msg interface{}, f func(ev cellnetEx.Event)) {
	if DefaultCell == nil {
		panic("RegitserModuleMsg Default nil")
	}
	DefaultCell.RegisterMessage(msg, f)
}

//RegisterMessage 注册消息回调
func (bcell *BaseCell) RegisterMessage(msg interface{}, f func(ev cellnetEx.Event)) {
	bcell.msgHandler[reflect.TypeOf(msg)] = f
}

//RegitserPlayerPBMessage 注册默认消息响应
func RegitserPlayerPBMessage(player interface{}) {
	if DefaultCell == nil {
		panic("RegitserModuleMsg Default nil")
	}
	DefaultCell.RegitserPlayerPBMessage(player)
}

//RegitserPlayerPBMessage 注册玩家处理的消息
func (bcell *BaseCell) RegitserPlayerPBMessage(player interface{}) {
	typeInfo := reflect.TypeOf(player)
	if typeInfo.Kind() != reflect.Ptr {
		panic("player must ptr")
	}
	if _, exsit := typeInfo.MethodByName("QID"); !exsit {
		panic("player must have QID Method")
	}
	if _, exsit := typeInfo.MethodByName("UID"); !exsit {
		panic("player must have UID Method")
	}
	for i := 0; i < typeInfo.NumMethod(); i++ {
		method := typeInfo.Method(i)
		if method.Type.NumIn() != 2 {
			continue
		}

		if cellnetEx.MessageMetaByType(method.Type.In(1)) == nil {
			continue
		}
		index := i
		msg := reflect.New(method.Type.In(1).Elem()).Interface()
		bcell.msgHandler[reflect.TypeOf(msg)] = func(ev cellnetEx.Event) {
			if ev.Session().GetUserData() == nil {
				log.Warnln("OnPlayerMessage not login close session", ev.Session().ID())
				ev.Session().Close()
				return
			}
			in := []reflect.Value{reflect.ValueOf(ev.Message())}
			reflect.ValueOf(ev.Session().GetUserData()).Method(index).Call(in)
		}
	}
}

//RegisterObjMessge 注册玩家相关的模块消息响应
func RegisterObjMessge(player interface{}) {
	if DefaultCell == nil {
		panic("RegitserModuleMsg Default nil")
	}
	DefaultCell.RegisterObjMessge(player)
}

//RegisterObjMessge 注册玩家下对象处理的消息
func (bcell *BaseCell) RegisterObjMessge(obj interface{}) {
	typeInfo := reflect.TypeOf(obj)
	if typeInfo.Kind() != reflect.Ptr {
		panic("obj must ptr")
	}
	for i := 0; i < typeInfo.NumMethod(); i++ {
		method := typeInfo.Method(i)
		if method.Type.NumIn() != 3 {
			continue
		}

		if method.Type.In(1).Kind() != reflect.String ||
			cellnetEx.MessageMetaByType(method.Type.In(2)) == nil {
			continue
		}

		index := i
		msg := reflect.New(method.Type.In(2).Elem()).Interface()
		bcell.msgHandler[reflect.TypeOf(msg)] = func(ev cellnetEx.Event) {
			userData := ev.Session().GetUserData()
			if userData == nil {
				log.Warnln("OnPlayerMessage Obj not login close session", ev.Session().ID())
				ev.Session().Close()
				return
			}
			in := []reflect.Value{
				reflect.ValueOf(userData.(iUserData).UID()),
				reflect.ValueOf(ev.Message()),
			}
			obj := reflect.ValueOf(userData).Elem().FieldByName(typeInfo.Elem().Name())
			obj.Method(index).Call(in)
		}
	}
}
