package basecell

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/davyxu/golog"

	"github.com/adamluo159/cellnetEx"
)

var log *golog.Logger = golog.New("basecell")

//DefaultCell 默认服务
var DefaultCell *BaseCell = nil

//IModule 模块接口
type IModule interface {
	Init()
	Name() string
	OnDestory()
}

//iUserData 用户数据接口
type IUserData interface {
	QID() int
	UID() string
}

//BaseCell 基础服务
type BaseCell struct {
	MsgQueueLen int
	modules     []IModule
	msgHandler  map[reflect.Type]func(ev cellnetEx.Event)
	mainQueue   cellnetEx.EventQueue
	msgQueues   []cellnetEx.EventQueue
	peer        cellnetEx.GenericPeer
	authCmdType reflect.Type //认证消息在客户端发的第一个消息
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
		mainQueue:   cellnetEx.NewEventQueue(),
		msgQueues:   make([]cellnetEx.EventQueue, 0),
		msgHandler:  make(map[reflect.Type]func(ev cellnetEx.Event)),
	}

	bcell.mainQueue.EnableCapturePanic(true)
	for i := 0; i < msgQueLen; i++ {
		q := cellnetEx.NewEventQueue()
		q.EnableCapturePanic(true)
		bcell.msgQueues = append(bcell.msgQueues, q)
	}

	if DefaultCell == nil {
		DefaultCell = bcell
	}
	return bcell
}

//初始化认证信息 应该是客户端发的第一个消息
func InitAuthMessage(authMessage interface{}) {
	if DefaultCell == nil {
		panic("RegitserModuleMsg Default nil")
	}
	DefaultCell.InitAuthMessage(authMessage)
}

//初始化认证信息 应该是客户端发的第一个消息
func (bcell *BaseCell) InitAuthMessage(authMessage interface{}) {
	bcell.authCmdType = reflect.TypeOf(authMessage)
}

func (bcell *BaseCell) msgQueue() func(ev cellnetEx.Event) {
	return func(ev cellnetEx.Event) {
		cmdType := reflect.TypeOf(ev.Message())
		qid := 0
		udata := ev.Session().GetUserData()
		if udata == nil {
			if cmdType != bcell.authCmdType {
				log.Warnf("frist Client Message should %s  current:%s", cmdType.String(), bcell.authCmdType.String())
				return
			}
			qid = rand.Intn(bcell.MsgQueueLen)
		} else {
			qid = udata.(IUserData).QID()
		}
		f, ok := bcell.msgHandler[cmdType]
		if !ok {
			log.Errorln("onMessage not found message handler ", ev.Message())
			return
		}
		bcell.msgQueues[qid].Post(func() {
			f(ev)
		})
	}
}

//Start 服务开始
func (bcell *BaseCell) Start(mods ...IModule) {
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
	if bcell.authCmdType == nil {
		panic("InitAuthMessage must set")
	}

	bcell.modules = mods
	// 开始侦听
	bcell.peer.Start()

	// 事件队列开始循环
	bcell.mainQueue.StartLoop()

	for _, v := range bcell.msgQueues {
		v.StartLoop()
	}
}

//Stop 服务停止
func (bcell *BaseCell) Stop() {
	bcell.peer.Stop()
	bcell.mainQueue.StopLoop()
	bcell.mainQueue.Wait()

	for _, v := range bcell.msgQueues {
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
		if method.Type.NumIn() != 2 {
			continue
		}

		if cellnetEx.MessageMetaByType(method.Type.In(1)) == nil {
			continue
		}

		index := i
		msgType := method.Type.In(1).Elem()
		msg := reflect.New(msgType).Interface()
		bcell.msgHandler[reflect.TypeOf(msg)] = func(ev cellnetEx.Event) {
			userData := ev.Session().GetUserData()
			if userData == nil {
				log.Warnln("RegisterObjMessge Obj not login close session", ev.Session().ID())
				ev.Session().Close()
				return
			}
			in := []reflect.Value{
				reflect.ValueOf(ev.Message()),
			}
			obj := reflect.ValueOf(userData).Elem().FieldByName(typeInfo.Elem().Name())
			if !obj.IsValid() {
				log.Errorf("RegisterObjMessge player field:%s not exsit drop message:%s", typeInfo.Elem().Name(), msgType.String())
				return
			}
			if obj.IsNil() {
				log.Errorf("RegisterObjMessge player field:%s nil drop message:%s", typeInfo.Elem().Name(), msgType.String())
				return
			}
			obj.Elem().Method(index).Call(in)
		}
	}
}
