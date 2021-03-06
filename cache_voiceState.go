package disgord

import (
	"errors"

	"github.com/andersfylling/disgord/cache/interfaces"
)

func createVoiceStateCacher(conf *CacheConfig) (cacher interfaces.CacheAlger, err error) {
	if conf.DisableVoiceStateCaching {
		return nil, nil
	}

	cacher, err = constructSpecificCacher(conf.VoiceStateCacheAlgorithm, 0, conf.VoiceStateCacheLifetime)
	return
}

type guildVoiceStatesCache struct {
	sessions []*VoiceState
}

func (g *guildVoiceStatesCache) existingSession(state *VoiceState) bool {
	return g.sessionPosition(state) > -1
}

func (g *guildVoiceStatesCache) sessionPosition(state *VoiceState) int {
	for i := range g.sessions {
		// if a channel is moved, the channelID should change(?)
		//if g.sessions[i].ChannelID != state.ChannelID {
		//	continue
		//}

		if g.sessions[i].UserID != state.UserID {
			continue
		}

		if g.sessions[i].SessionID == state.SessionID {
			return i
		}
	}

	return -1
}

func (g *guildVoiceStatesCache) update(state *VoiceState, copyOnWrite bool) {
	pos := g.sessionPosition(state)
	if state.ChannelID.Empty() {
		// someone left
		if pos > -1 {
			g.sessions[pos] = g.sessions[len(g.sessions)-1]
			g.sessions[len(g.sessions)-1] = nil
			g.sessions = g.sessions[:len(g.sessions)-1]
		}
		return
	}

	// someone joined / moved channel. But there exist no information about the session
	// so we register the information
	if pos < 0 {
		var data *VoiceState
		if copyOnWrite {
			data = state.DeepCopy().(*VoiceState) // TODO: handle member
		} else {
			data = state
		}
		g.sessions = append(g.sessions, data)
		return
	}

	// someone moved an existing channel
	if g.sessions[pos].ChannelID != state.ChannelID {
		g.sessions[pos].ChannelID = state.ChannelID
		return
	}

	// TODO: this point should not be reached, unless the above checks are incomplete
}

// SetVoiceState adds a new voice state to cache or updates an existing one
func (c *Cache) SetVoiceState(state *VoiceState) {
	if c.voiceStates == nil || state == nil {
		return
	}

	c.voiceStates.Lock()
	defer c.voiceStates.Unlock()

	id := state.GuildID
	if item, exists := c.voiceStates.Get(id); exists {
		states := item.Object().(*guildVoiceStatesCache)
		states.update(state, c.immutable)
		c.users.RefreshAfterDiscordUpdate(item)
	} else {
		states := &guildVoiceStatesCache{}
		states.update(state, c.immutable)
		c.voiceStates.Set(id, c.voiceStates.CreateCacheableItem(states))
	}
}

type guildVoiceStateCacheParams struct {
	userID    Snowflake
	channelID Snowflake
	sessionID string
}

// GetVoiceState ...
func (c *Cache) GetVoiceState(guildID Snowflake, params *guildVoiceStateCacheParams) (state *VoiceState, err error) {
	if c.voiceStates == nil {
		err = newErrorUsingDeactivatedCache("voice-states")
		return
	}

	c.voiceStates.RLock()
	defer c.voiceStates.RUnlock()

	var exists bool
	var result interfaces.CacheableItem
	if result, exists = c.voiceStates.Get(guildID); !exists {
		err = newErrorCacheItemNotFound(guildID)
		return
	}

	states := result.Object().(*guildVoiceStatesCache)
	filter := &VoiceState{
		ChannelID: params.channelID,
		UserID:    params.userID,
		SessionID: params.sessionID,
	}
	pos := states.sessionPosition(filter)
	if pos < 0 {
		err = errors.New("unable to find state with given params filter")
		return
	}

	match := states.sessions[pos]
	if c.immutable {
		state = match.DeepCopy().(*VoiceState)
	} else {
		state = match
	}

	return
}
