package crouter

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// PermToString turns a discordgo permission constant into a string
func PermToString(p int) string {
	switch p {
	case discordgo.PermissionSendMessages:
		return "Send Messages"
	case discordgo.PermissionSendTTSMessages:
		return "Send TTS Messages"
	case discordgo.PermissionManageMessages:
		return "Manage Messages"
	case discordgo.PermissionEmbedLinks:
		return "Embed Links"
	case discordgo.PermissionAttachFiles:
		return "Attach Files"
	case discordgo.PermissionReadMessageHistory:
		return "Read Message History"
	case discordgo.PermissionMentionEveryone:
		return "Mention @everyone, @here, and All Roles"
	case discordgo.PermissionUseExternalEmojis:
		return "Use External Emojis"
	case discordgo.PermissionVoiceConnect:
		return "Voice Connect"
	case discordgo.PermissionVoiceSpeak:
		return "Voice Speak"
	case discordgo.PermissionVoiceMuteMembers:
		return "Voice Mute Members"
	case discordgo.PermissionVoiceDeafenMembers:
		return "Voice Deafen Members"
	case discordgo.PermissionVoiceMoveMembers:
		return "Voice Move Members"
	case discordgo.PermissionVoiceUseVAD:
		return "Use Voice Activity"
	case discordgo.PermissionVoicePrioritySpeaker:
		return "Priority Speaker"
	case discordgo.PermissionChangeNickname:
		return "Change Nickname"
	case discordgo.PermissionManageNicknames:
		return "Manage Nicknames"
	case discordgo.PermissionManageRoles:
		return "Manage Roles"
	case discordgo.PermissionManageWebhooks:
		return "Manage Webhooks"
	case discordgo.PermissionManageEmojis:
		return "Manage Emojis"
	case discordgo.PermissionCreateInstantInvite:
		return "Create Invite"
	case discordgo.PermissionKickMembers:
		return "Kick Members"
	case discordgo.PermissionBanMembers:
		return "Ban Members"
	case discordgo.PermissionAdministrator:
		return "Administrator"
	case discordgo.PermissionManageChannels:
		return "Manage Channels"
	case discordgo.PermissionManageServer:
		return "Manage Server"
	case discordgo.PermissionAddReactions:
		return "Add Reactions"
	case discordgo.PermissionViewAuditLogs:
		return "View Audit Log"
	case discordgo.PermissionViewChannel:
		return "Read Text Channels & See Voice Channels"
	default:
		return fmt.Sprint(p)
	}
}
