package disgord

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	. "github.com/andersfylling/disgord/event"
)

// NewDispatch construct a Dispatch object for reacting to web socket events
// from discord
func NewDispatch() *Dispatch {
	dispatcher := &Dispatch{
		allChan:                      make(chan interface{}),
		readyChan:                    make(chan *Ready),
		resumedChan:                  make(chan *Resumed),
		channelCreateChan:            make(chan *ChannelCreate),
		channelUpdateChan:            make(chan *ChannelUpdate),
		channelDeleteChan:            make(chan *ChannelDelete),
		channelPinsUpdateChan:        make(chan *ChannelPinsUpdate),
		guildCreateChan:              make(chan *GuildCreate),
		guildUpdateChan:              make(chan *GuildUpdate),
		guildDeleteChan:              make(chan *GuildDelete),
		guildBanAddChan:              make(chan *GuildBanAdd),
		guildBanRemoveChan:           make(chan *GuildBanRemove),
		guildEmojisUpdateChan:        make(chan *GuildEmojisUpdate),
		guildIntegrationsUpdateChan:  make(chan *GuildIntegrationsUpdate),
		guildMemberAddChan:           make(chan *GuildMemberAdd),
		guildMemberRemoveChan:        make(chan *GuildMemberRemove),
		guildMemberUpdateChan:        make(chan *GuildMemberUpdate),
		guildMembersChunkChan:        make(chan *GuildMembersChunk),
		guildRoleUpdateChan:          make(chan *GuildRoleUpdate),
		guildRoleCreateChan:          make(chan *GuildRoleCreate),
		guildRoleDeleteChan:          make(chan *GuildRoleDelete),
		messageCreateChan:            make(chan *MessageCreate),
		messageUpdateChan:            make(chan *MessageUpdate),
		messageDeleteChan:            make(chan *MessageDelete),
		messageDeleteBulkChan:        make(chan *MessageDeleteBulk),
		messageReactionAddChan:       make(chan *MessageReactionAdd),
		messageReactionRemoveChan:    make(chan *MessageReactionRemove),
		messageReactionRemoveAllChan: make(chan *MessageReactionRemoveAll),
		presenceUpdateChan:           make(chan *PresenceUpdate),
		typingStartChan:              make(chan *TypingStart),
		userUpdateChan:               make(chan *UserUpdate),
		voiceStateUpdateChan:         make(chan *VoiceStateUpdate),
		voiceServerUpdateChan:        make(chan *VoiceServerUpdate),
		webhooksUpdateChan:           make(chan *WebhooksUpdate),

		listeners:      make(map[string][]interface{}),
		listenOnceOnly: make(map[string][]int),

		shutdown: make(chan struct{}),
	}

	return dispatcher
}

// EvtDispatcher interface for developers using the module. Gives access to
// channels and the option to add callbacks to specific event types.
type EvtDispatcher interface {
	AllChan() <-chan interface{} // any event
	ReadyChan() <-chan *Ready
	ResumedChan() <-chan *Resumed
	ChannelCreateChan() <-chan *ChannelCreate
	ChannelUpdateChan() <-chan *ChannelUpdate
	ChannelDeleteChan() <-chan *ChannelDelete
	ChannelPinsUpdateChan() <-chan *ChannelPinsUpdate
	GuildCreateChan() <-chan *GuildCreate
	GuildUpdateChan() <-chan *GuildUpdate
	GuildDeleteChan() <-chan *GuildDelete
	GuildBanAddChan() <-chan *GuildBanAdd
	GuildBanRemoveChan() <-chan *GuildBanRemove
	GuildEmojisUpdateChan() <-chan *GuildEmojisUpdate
	GuildIntegrationsUpdateChan() <-chan *GuildIntegrationsUpdate
	GuildMemberAddChan() <-chan *GuildMemberAdd
	GuildMemberRemoveChan() <-chan *GuildMemberRemove
	GuildMemberUpdateChan() <-chan *GuildMemberUpdate
	GuildMembersChunkChan() <-chan *GuildMembersChunk
	GuildRoleUpdateChan() <-chan *GuildRoleUpdate
	GuildRoleCreateChan() <-chan *GuildRoleCreate
	GuildRoleDeleteChan() <-chan *GuildRoleDelete
	MessageCreateChan() <-chan *MessageCreate
	MessageUpdateChan() <-chan *MessageUpdate
	MessageDeleteChan() <-chan *MessageDelete
	MessageDeleteBulkChan() <-chan *MessageDeleteBulk
	MessageReactionAddChan() <-chan *MessageReactionAdd
	MessageReactionRemoveChan() <-chan *MessageReactionRemove
	MessageReactionRemoveAllChan() <-chan *MessageReactionRemoveAll
	PresenceUpdateChan() <-chan *PresenceUpdate
	TypingStartChan() <-chan *TypingStart
	UserUpdateChan() <-chan *UserUpdate
	VoiceStateUpdateChan() <-chan *VoiceStateUpdate
	VoiceServerUpdateChan() <-chan *VoiceServerUpdate
	WebhooksUpdateChan() <-chan *WebhooksUpdate

	AddHandler(evtName string, listener interface{})
	AddHandlerOnce(evtName string, listener interface{})
}

// Dispatch holds all the channels and internal state for all registered
// observers
type Dispatch struct {
	allChan                      chan interface{} // any event
	readyChan                    chan *Ready
	resumedChan                  chan *Resumed
	channelCreateChan            chan *ChannelCreate
	channelUpdateChan            chan *ChannelUpdate
	channelDeleteChan            chan *ChannelDelete
	channelPinsUpdateChan        chan *ChannelPinsUpdate
	guildCreateChan              chan *GuildCreate
	guildUpdateChan              chan *GuildUpdate
	guildDeleteChan              chan *GuildDelete
	guildBanAddChan              chan *GuildBanAdd
	guildBanRemoveChan           chan *GuildBanRemove
	guildEmojisUpdateChan        chan *GuildEmojisUpdate
	guildIntegrationsUpdateChan  chan *GuildIntegrationsUpdate
	guildMemberAddChan           chan *GuildMemberAdd
	guildMemberRemoveChan        chan *GuildMemberRemove
	guildMemberUpdateChan        chan *GuildMemberUpdate
	guildMembersChunkChan        chan *GuildMembersChunk
	guildRoleUpdateChan          chan *GuildRoleUpdate
	guildRoleCreateChan          chan *GuildRoleCreate
	guildRoleDeleteChan          chan *GuildRoleDelete
	messageCreateChan            chan *MessageCreate
	messageUpdateChan            chan *MessageUpdate
	messageDeleteChan            chan *MessageDelete
	messageDeleteBulkChan        chan *MessageDeleteBulk
	messageReactionAddChan       chan *MessageReactionAdd
	messageReactionRemoveChan    chan *MessageReactionRemove
	messageReactionRemoveAllChan chan *MessageReactionRemoveAll
	presenceUpdateChan           chan *PresenceUpdate
	typingStartChan              chan *TypingStart
	userUpdateChan               chan *UserUpdate
	voiceStateUpdateChan         chan *VoiceStateUpdate
	voiceServerUpdateChan        chan *VoiceServerUpdate
	webhooksUpdateChan           chan *WebhooksUpdate

	listeners      map[string][]interface{}
	listenOnceOnly map[string][]int

	shutdown chan struct{}

	listenersLock sync.RWMutex
}

func (d *Dispatch) start() {
	// make sure every channel has a receiver to avoid deadlock
	// TODO: review, this feels hacky
	d.alwaysListenToChans()
}

func (d *Dispatch) stop() {
	close(d.shutdown)
}

// On places listeners into their respected stacks
// func (d *Dispatcher) OnEvent(evtName string, listener EventCallback) {
// 	d.listeners[evtName] = append(d.listeners[evtName], listener)
// }

// alwaysListenToChans makes sure no deadlocks occure
func (d *Dispatch) alwaysListenToChans() {
	go func() {
		stop := false
		for {
			select {
			case <-d.allChan:
			case <-d.readyChan:
			case <-d.resumedChan:
			case <-d.channelCreateChan:
			case <-d.channelDeleteChan:
			case <-d.channelPinsUpdateChan:
			case <-d.channelUpdateChan:
			case <-d.guildBanAddChan:
			case <-d.guildBanRemoveChan:
			case <-d.guildCreateChan:
			case <-d.guildDeleteChan:
			case <-d.guildEmojisUpdateChan:
			case <-d.guildIntegrationsUpdateChan:
			case <-d.guildMemberAddChan:
			case <-d.guildMemberRemoveChan:
			case <-d.guildMemberUpdateChan:
			case <-d.guildMembersChunkChan:
			case <-d.guildRoleCreateChan:
			case <-d.guildRoleDeleteChan:
			case <-d.guildRoleUpdateChan:
			case <-d.guildUpdateChan:
			case <-d.messageCreateChan:
			case <-d.messageDeleteBulkChan:
			case <-d.messageDeleteChan:
			case <-d.messageReactionAddChan:
			case <-d.messageReactionRemoveAllChan:
			case <-d.messageReactionRemoveChan:
			case <-d.messageUpdateChan:
			case <-d.presenceUpdateChan:
			case <-d.typingStartChan:
			case <-d.userUpdateChan:
			case <-d.voiceStateUpdateChan:
			case <-d.voiceServerUpdateChan:
			case <-d.webhooksUpdateChan:
			case <-d.shutdown:
				stop = true
			}

			if stop {
				break
			}
		}
	}()
}

func (d *Dispatch) triggerChan(ctx context.Context, evtName string, session Session, box interface{}) {
	switch evtName {
	case KeyReady:
		d.readyChan <- box.(*Ready)
	case KeyResumed:
		d.resumedChan <- box.(*Resumed)
	case KeyChannelCreate:
		d.channelCreateChan <- box.(*ChannelCreate)
	case KeyChannelUpdate:
		d.channelUpdateChan <- box.(*ChannelUpdate)
	case KeyChannelDelete:
		d.channelDeleteChan <- box.(*ChannelDelete)
	case KeyChannelPinsUpdate:
		d.channelPinsUpdateChan <- box.(*ChannelPinsUpdate)
	case KeyGuildCreate:
		d.guildCreateChan <- box.(*GuildCreate)
	case KeyGuildUpdate:
		d.guildUpdateChan <- box.(*GuildUpdate)
	case KeyGuildDelete:
		d.guildDeleteChan <- box.(*GuildDelete)
	case KeyGuildBanAdd:
		d.guildBanAddChan <- box.(*GuildBanAdd)
	case KeyGuildBanRemove:
		d.guildBanRemoveChan <- box.(*GuildBanRemove)
	case KeyGuildEmojisUpdate:
		d.guildEmojisUpdateChan <- box.(*GuildEmojisUpdate)
	case KeyGuildIntegrationsUpdate:
		d.guildIntegrationsUpdateChan <- box.(*GuildIntegrationsUpdate)
	case KeyGuildMemberAdd:
		d.guildMemberAddChan <- box.(*GuildMemberAdd)
	case KeyGuildMemberRemove:
		d.guildMemberRemoveChan <- box.(*GuildMemberRemove)
	case KeyGuildMemberUpdate:
		d.guildMemberUpdateChan <- box.(*GuildMemberUpdate)
	case KeyGuildMembersChunk:
		d.guildMembersChunkChan <- box.(*GuildMembersChunk)
	case KeyGuildRoleCreate:
		d.guildRoleCreateChan <- box.(*GuildRoleCreate)
	case KeyGuildRoleUpdate:
		d.guildRoleUpdateChan <- box.(*GuildRoleUpdate)
	case KeyGuildRoleDelete:
		d.guildRoleDeleteChan <- box.(*GuildRoleDelete)
	case KeyMessageCreate:
		d.messageCreateChan <- box.(*MessageCreate)
	case KeyMessageUpdate:
		d.messageUpdateChan <- box.(*MessageUpdate)
	case KeyMessageDelete:
		d.messageDeleteChan <- box.(*MessageDelete)
	case KeyMessageDeleteBulk:
		d.messageDeleteBulkChan <- box.(*MessageDeleteBulk)
	case KeyMessageReactionAdd:
		d.messageReactionAddChan <- box.(*MessageReactionAdd)
	case KeyMessageReactionRemove:
		d.messageReactionRemoveChan <- box.(*MessageReactionRemove)
	case KeyMessageReactionRemoveAll:
		d.messageReactionRemoveAllChan <- box.(*MessageReactionRemoveAll)
	case KeyPresenceUpdate:
		d.presenceUpdateChan <- box.(*PresenceUpdate)
	case KeyTypingStart:
		d.typingStartChan <- box.(*TypingStart)
	case KeyUserUpdate:
		d.userUpdateChan <- box.(*UserUpdate)
	case KeyVoiceStateUpdate:
		d.voiceStateUpdateChan <- box.(*VoiceStateUpdate)
	case KeyVoiceServerUpdate:
		d.voiceServerUpdateChan <- box.(*VoiceServerUpdate)
	case KeyWebhooksUpdate:
		d.webhooksUpdateChan <- box.(*WebhooksUpdate)
	default:
		fmt.Printf("------\nTODO\nImplement channel for `%s`\n------\n\n", evtName)
	}
}

func (d *Dispatch) triggerCallbacks(ctx context.Context, evtName string, session Session, box interface{}) {
	switch evtName {
	case KeyReady:
		for _, listener := range d.listeners[KeyReady] {
			go (listener.(ReadyCallback))(session, box.(*Ready))
		}
	case KeyResumed:
		for _, listener := range d.listeners[KeyResumed] {
			go (listener.(ResumedCallback))(session, box.(*Resumed))
		}
	case KeyChannelCreate:
		for _, listener := range d.listeners[KeyChannelCreate] {
			go (listener.(ChannelCreateCallback))(session, box.(*ChannelCreate))
		}
	case KeyChannelUpdate:
		for _, listener := range d.listeners[KeyChannelUpdate] {
			go (listener.(ChannelUpdateCallback))(session, box.(*ChannelUpdate))
		}
	case KeyChannelDelete:
		for _, listener := range d.listeners[KeyChannelDelete] {
			go (listener.(ChannelDeleteCallback))(session, box.(*ChannelDelete))
		}
	case KeyChannelPinsUpdate:
		for _, listener := range d.listeners[KeyChannelPinsUpdate] {
			go (listener.(ChannelPinsUpdateCallback))(session, box.(*ChannelPinsUpdate))
		}
	case KeyGuildCreate:
		for _, listener := range d.listeners[KeyGuildCreate] {
			go (listener.(GuildCreateCallback))(session, box.(*GuildCreate))
		}
	case KeyGuildUpdate:
		for _, listener := range d.listeners[KeyGuildUpdate] {
			go (listener.(GuildUpdateCallback))(session, box.(*GuildUpdate))
		}
	case KeyGuildDelete:
		for _, listener := range d.listeners[KeyGuildDelete] {
			go (listener.(GuildDeleteCallback))(session, box.(*GuildDelete))
		}
	case KeyGuildBanAdd:
		for _, listener := range d.listeners[KeyGuildBanAdd] {
			go (listener.(GuildBanAddCallback))(session, box.(*GuildBanAdd))
		}
	case KeyGuildBanRemove:
		for _, listener := range d.listeners[KeyGuildBanRemove] {
			go (listener.(GuildBanRemoveCallback))(session, box.(*GuildBanRemove))
		}
	case KeyGuildEmojisUpdate:
		for _, listener := range d.listeners[KeyGuildEmojisUpdate] {
			go (listener.(GuildEmojisUpdateCallback))(session, box.(*GuildEmojisUpdate))
		}
	case KeyGuildIntegrationsUpdate:
		for _, listener := range d.listeners[KeyGuildIntegrationsUpdate] {
			go (listener.(GuildIntegrationsUpdateCallback))(session, box.(*GuildIntegrationsUpdate))
		}
	case KeyGuildMemberAdd:
		for _, listener := range d.listeners[KeyGuildMemberAdd] {
			go (listener.(GuildMemberAddCallback))(session, box.(*GuildMemberAdd))
		}
	case KeyGuildMemberRemove:
		for _, listener := range d.listeners[KeyGuildMemberRemove] {
			go (listener.(GuildMemberRemoveCallback))(session, box.(*GuildMemberRemove))
		}
	case KeyGuildMemberUpdate:
		for _, listener := range d.listeners[KeyGuildMemberUpdate] {
			go (listener.(GuildMemberUpdateCallback))(session, box.(*GuildMemberUpdate))
		}
	case KeyGuildMembersChunk:
		for _, listener := range d.listeners[KeyGuildMembersChunk] {
			go (listener.(GuildMembersChunkCallback))(session, box.(*GuildMembersChunk))
		}
	case KeyGuildRoleCreate:
		for _, listener := range d.listeners[KeyGuildRoleCreate] {
			go (listener.(GuildRoleCreateCallback))(session, box.(*GuildRoleCreate))
		}
	case KeyGuildRoleUpdate:
		for _, listener := range d.listeners[KeyGuildRoleUpdate] {
			go (listener.(GuildRoleUpdateCallback))(session, box.(*GuildRoleUpdate))
		}
	case KeyGuildRoleDelete:
		for _, listener := range d.listeners[KeyGuildRoleDelete] {
			go (listener.(GuildRoleDeleteCallback))(session, box.(*GuildRoleDelete))
		}
	case KeyMessageCreate:
		for _, listener := range d.listeners[KeyMessageCreate] {
			go (listener.(MessageCreateCallback))(session, box.(*MessageCreate))
		}
	case KeyMessageUpdate:
		for _, listener := range d.listeners[KeyMessageUpdate] {
			go (listener.(MessageUpdateCallback))(session, box.(*MessageUpdate))
		}
	case KeyMessageDelete:
		for _, listener := range d.listeners[KeyMessageDelete] {
			go (listener.(MessageDeleteCallback))(session, box.(*MessageDelete))
		}
	case KeyMessageDeleteBulk:
		for _, listener := range d.listeners[KeyMessageDeleteBulk] {
			go (listener.(MessageDeleteBulkCallback))(session, box.(*MessageDeleteBulk))
		}
	case KeyMessageReactionAdd:
		for _, listener := range d.listeners[KeyMessageReactionAdd] {
			go (listener.(MessageReactionAddCallback))(session, box.(*MessageReactionAdd))
		}
	case KeyMessageReactionRemove:
		for _, listener := range d.listeners[KeyMessageReactionRemove] {
			go (listener.(MessageReactionRemoveCallback))(session, box.(*MessageReactionRemove))
		}
	case KeyMessageReactionRemoveAll:
		for _, listener := range d.listeners[KeyMessageReactionRemoveAll] {
			go (listener.(MessageReactionRemoveAllCallback))(session, box.(*MessageReactionRemoveAll))
		}
	case KeyPresenceUpdate:
		for _, listener := range d.listeners[KeyPresenceUpdate] {
			go (listener.(PresenceUpdateCallback))(session, box.(*PresenceUpdate))
		}
	case KeyTypingStart:
		for _, listener := range d.listeners[KeyTypingStart] {
			go (listener.(TypingStartCallback))(session, box.(*TypingStart))
		}
	case KeyUserUpdate:
		for _, listener := range d.listeners[KeyUserUpdate] {
			go (listener.(UserUpdateCallback))(session, box.(*UserUpdate))
		}
	case KeyVoiceStateUpdate:
		for _, listener := range d.listeners[KeyVoiceStateUpdate] {
			go (listener.(VoiceStateUpdateCallback))(session, box.(*VoiceStateUpdate))
		}
	case KeyVoiceServerUpdate:
		for _, listener := range d.listeners[KeyVoiceServerUpdate] {
			go (listener.(VoiceServerUpdateCallback))(session, box.(*VoiceServerUpdate))
		}
	case KeyWebhooksUpdate:
		for _, listener := range d.listeners[KeyWebhooksUpdate] {
			go (listener.(WebhooksUpdateCallback))(session, box.(*WebhooksUpdate))
		}
	default:
		fmt.Printf("------\nTODO\nImplement callback for `%s`\n------\n\n", evtName)
	}

	// remove the run only once listeners
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	for _, index := range d.listenOnceOnly[evtName] {
		// https://github.com/golang/go/wiki/SliceTricks#delete-without-preserving-order
		d.listeners[evtName][index] = d.listeners[evtName][len(d.listeners[evtName])-1]
		d.listeners[evtName][len(d.listeners[evtName])-1] = nil
		d.listeners[evtName] = d.listeners[evtName][:len(d.listeners[evtName])-1]
	}

	// remove the once only register
	_, exists := d.listenOnceOnly[evtName]
	if exists {
		delete(d.listenOnceOnly, evtName)
	}
}

// AllChan sends all event types
func (d *Dispatch) AllChan() <-chan interface{} {
	return d.allChan
}

// ReadyChan for READY events
func (d *Dispatch) ReadyChan() <-chan *Ready {
	return d.readyChan
}

// ResumedChan for RESUME events
func (d *Dispatch) ResumedChan() <-chan *Resumed {
	return d.resumedChan
}

// ChannelCreateChan for CHANNEL_CREATE, when a channel is created
func (d *Dispatch) ChannelCreateChan() <-chan *ChannelCreate {
	return d.channelCreateChan
}

// ChannelUpdateChan for CHANNEL_UPDATE, when a channel is updated
func (d *Dispatch) ChannelUpdateChan() <-chan *ChannelUpdate {
	return d.channelUpdateChan
}

// ChannelDeleteChan for CHANNEL_DELETE, when a channel is deleted
func (d *Dispatch) ChannelDeleteChan() <-chan *ChannelDelete {
	return d.channelDeleteChan
}

// ChannelPinsUpdateChan for CHANNEL_PIN_UPDATE. Message was pinned or unpinned
func (d *Dispatch) ChannelPinsUpdateChan() <-chan *ChannelPinsUpdate {
	return d.channelPinsUpdateChan
}

// GuildCreateChan for GUILD_CREATE. Lazy-load for unavailable guild, guild
// became available, or user joined a new guild
func (d *Dispatch) GuildCreateChan() <-chan *GuildCreate {
	return d.guildCreateChan
}

// GuildUpdateChan for GUILD_UPDATE. Guild was updated
func (d *Dispatch) GuildUpdateChan() <-chan *GuildUpdate {
	return d.guildUpdateChan
}

// GuildDeleteChan for GUILD_DELETE, guild became unavailable, or user
// left/was removed from a guild
func (d *Dispatch) GuildDeleteChan() <-chan *GuildDelete {
	return d.guildDeleteChan
}

// GuildBanAddChan for GUILD_BAN_ADD. A user was banned from a guild
func (d *Dispatch) GuildBanAddChan() <-chan *GuildBanAdd {
	return d.guildBanAddChan
}

// GuildBanRemoveChan for GUILD_BAN_REMOVE. A user was unbanned from a guild
func (d *Dispatch) GuildBanRemoveChan() <-chan *GuildBanRemove {
	return d.guildBanRemoveChan
}

// GuildEmojisUpdateChan for GUILD_EMOJI_UPDATE. Guild emojis were updated
func (d *Dispatch) GuildEmojisUpdateChan() <-chan *GuildEmojisUpdate {
	return d.guildEmojisUpdateChan
}

// GuildIntegrationsUpdateChan for GUILD_INTEGRATIONS_UPDATE. Guild integration
// was updated
func (d *Dispatch) GuildIntegrationsUpdateChan() <-chan *GuildIntegrationsUpdate {
	return d.guildIntegrationsUpdateChan
}

// GuildMemberAddChan for GUILD_MEMBER_ADD. New user joined a guild.
func (d *Dispatch) GuildMemberAddChan() <-chan *GuildMemberAdd {
	return d.guildMemberAddChan
}

// GuildMemberRemoveChan for GUILD_MEMBER_REMOVE. User was removed from guild.
func (d *Dispatch) GuildMemberRemoveChan() <-chan *GuildMemberRemove {
	return d.guildMemberRemoveChan
}

// GuildMemberUpdateChan for GUILD_MEMBER_UPDATE. Guild member was updated.
func (d *Dispatch) GuildMemberUpdateChan() <-chan *GuildMemberUpdate {
	return d.guildMemberUpdateChan
}

// GuildMembersChunkChan for GUILD_MEMBERS_CHUNK. Response to socket command
// 'Request Guild Members'
func (d *Dispatch) GuildMembersChunkChan() <-chan *GuildMembersChunk {
	return d.guildMembersChunkChan
}

// GuildRoleCreateChan for GUILD_ROLE_CREATE. Guild role was created.
func (d *Dispatch) GuildRoleCreateChan() <-chan *GuildRoleCreate {
	return d.guildRoleCreateChan
}

// GuildRoleUpdateChan for GUILD_ROLE_UPDATE. Guild role was updated.
func (d *Dispatch) GuildRoleUpdateChan() <-chan *GuildRoleUpdate {
	return d.guildRoleUpdateChan
}

// GuildRoleDeleteChan for GUILD_ROLE_DELETE. Guild role was deleted.
func (d *Dispatch) GuildRoleDeleteChan() <-chan *GuildRoleDelete {
	return d.guildRoleDeleteChan
}

// MessageCreateChan for MESSAGE_CREATE. New message was created.
func (d *Dispatch) MessageCreateChan() <-chan *MessageCreate {
	return d.messageCreateChan
}

// MessageUpdateChan for MESSAGE_UPDATE. Message was updated.
func (d *Dispatch) MessageUpdateChan() <-chan *MessageUpdate {
	return d.messageUpdateChan
}

// MessageDeleteChan for MESSAGE_DELETE. Message was deleted.
func (d *Dispatch) MessageDeleteChan() <-chan *MessageDelete {
	return d.messageDeleteChan
}

// MessageDeleteBulkChan for MESSAGE_DELETE_BULK. Multiple messages were
// deleted at once.
func (d *Dispatch) MessageDeleteBulkChan() <-chan *MessageDeleteBulk {
	return d.messageDeleteBulkChan
}

// MessageReactionAddChan for MESSAGE_REACTION_ADD. A user reacted to a message.
func (d *Dispatch) MessageReactionAddChan() <-chan *MessageReactionAdd {
	return d.messageReactionAddChan
}

// MessageReactionRemoveChan for MESSAGE_REACTION_REMOVE. A user removed a
// a reaction to a message.
func (d *Dispatch) MessageReactionRemoveChan() <-chan *MessageReactionRemove {
	return d.messageReactionRemoveChan
}

// MessageReactionRemoveAllChan for MESSAGE_REACTION_REMOVE_ALL. All reactions
// were explicitly removed from a message
func (d *Dispatch) MessageReactionRemoveAllChan() <-chan *MessageReactionRemoveAll {
	return d.messageReactionRemoveAllChan
}

// PresenceUpdateChan for PRESENCE_UPDATE. A user's presence was updated in a
// guild.
func (d *Dispatch) PresenceUpdateChan() <-chan *PresenceUpdate {
	return d.presenceUpdateChan
}

// TypingStartChan for TYPING_START. A user started typing in a channel.
func (d *Dispatch) TypingStartChan() <-chan *TypingStart {
	return d.typingStartChan
}

// UserUpdateChan for USER_UPDATE. Properties about a user changed
func (d *Dispatch) UserUpdateChan() <-chan *UserUpdate {
	return d.userUpdateChan
}

// VoiceStateUpdateChan for VOICE_STATE_UPDATE. Someone joined, left, or moved
// a voice channel
func (d *Dispatch) VoiceStateUpdateChan() <-chan *VoiceStateUpdate {
	return d.voiceStateUpdateChan
}

// VoiceServerUpdateChan for VOICE_SERVER_UPDATE. Guild's voice server was
// updated
func (d *Dispatch) VoiceServerUpdateChan() <-chan *VoiceServerUpdate {
	return d.voiceServerUpdateChan
}

// WebhooksUpdateChan for WEBHOOK_UPDATE. A guild channel webhook was created,
// update, or deleted
func (d *Dispatch) WebhooksUpdateChan() <-chan *WebhooksUpdate {
	return d.webhooksUpdateChan
}

// AddHandler bind a callback to a event type. See event.Key...
func (d *Dispatch) AddHandler(evtName string, listener interface{}) {
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	d.listeners[evtName] = append(d.listeners[evtName], listener)
}

// AddHandlerOnce same as AddHandler except that the callback is only
// triggered once
func (d *Dispatch) AddHandlerOnce(evtName string, listener interface{}) {
	d.listenersLock.Lock()
	defer d.listenersLock.Unlock()

	index := len(d.listeners[evtName])
	d.listeners[evtName] = append(d.listeners[evtName], listener)
	d.listenOnceOnly[evtName] = append(d.listenOnceOnly[evtName], index)
}

// Unmarshal ...
// TODO: review
func Unmarshal(data []byte, box interface{}) {
	err := json.Unmarshal(data, box)
	if err != nil {
		panic(err) // !
	}
}
