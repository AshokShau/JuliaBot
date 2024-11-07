package modules

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
)

var startTime = time.Now()

func StartHandle(m *telegram.NewMessage) error {
	m.Reply("Hellow!")
	_, err := m.Client.MessagesSendReaction(&telegram.MessagesSendReactionParams{
		Peer:  m.Peer,
		Big:   true,
		MsgID: m.ID,
		Reaction: []telegram.Reaction{
			&telegram.ReactionEmoji{
				Emoticon: "❤️",
			},
		},
	})

	return err
}

func GatherSystemInfo(m *telegram.NewMessage) error {
	pid := os.Getpid()
	os := "<b>---- System Data ----</b>\n\n<b>Uptime:</b> " + time.Since(startTime).String() + "\n<b>PowerOS:</b> " + runtime.GOOS + "\n"
	os += "<b>Arch:</b> " + runtime.GOARCH + "\n"
	os += fmt.Sprintf("<b>CPUs:</b> %d\n", runtime.NumCPU())
	os += fmt.Sprintf("<b>Goroutines:</b> %d\n", runtime.NumGoroutine())
	cmd_to_get_currentprocess_memory := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "rss=")
	output, err := cmd_to_get_currentprocess_memory.Output()

	if err != nil {
		os += "<b>Error:</b> " + err.Error() + "\n"
	} else {
		os += "<b>Memory Usage:</b> " + strings.TrimSpace(string(output)) + " KB\n"
	}

	cmd_to_get_currentprocess_cpu := exec.Command("ps", "-p", strconv.Itoa(pid), "-o", "%cpu=")
	output, err = cmd_to_get_currentprocess_cpu.Output()

	if err != nil {
		os += "<b>Error:</b> " + err.Error() + "\n"
	} else {
		os += "<b>CPU Usage:</b> " + strings.TrimSpace(string(output)) + " %\n"
	}

	_, err = m.Reply(os)
	return err
}

func UserHandle(m *telegram.NewMessage) error {
	var userID int64 = 0
	var userHash int64 = 0
	if m.IsReply() {
		r, _ := m.GetReplyMessage()
		userID = r.SenderID()
		if r.Sender == nil {
			m.Reply("Error: User not found")
			return nil
		}
		userHash = r.Sender.AccessHash
	} else if len(m.Args()) > 0 {
		i, ok := strconv.Atoi(m.Args())
		if ok != nil {
			user, err := m.Client.ResolveUsername(m.Args())
			if err != nil {
				m.Reply("Error: " + err.Error())
				return nil
			}
			ux, ok := user.(*telegram.UserObj)
			if !ok {
				m.Reply("Error: User not found")
				return nil
			}
			userID = ux.ID
			userHash = ux.AccessHash
		} else {
			userID = int64(i)
			user, err := m.Client.GetUser(int64(i))
			if err != nil {
				m.Reply("Error: " + err.Error())
				return nil
			}
			userHash = user.AccessHash
		}
	} else {
		userID = m.SenderID()
		if m.Sender == nil {
			m.Reply("Error: User not found")
			return nil
		}
		userHash = m.Sender.AccessHash
	}
	user, err := m.Client.UsersGetFullUser(&telegram.InputUserObj{
		UserID:     userID,
		AccessHash: userHash,
	})

	if err != nil {
		m.Reply("Error: " + err.Error())
		return nil
	}

	uf := user.FullUser
	un := user.Users[0].(*telegram.UserObj)

	var userString string
	userString += "<b>---- User Info ----</b>\n\n"
	if un.FirstName != "" {
		userString += "<b>First Name:</b> " + un.FirstName + "\n"
	}
	if un.LastName != "" {
		userString += "<b>Last Name:</b> " + un.LastName + "\n"
	}
	if un.Username != "" {
		userString += "<b>Username:</b> @" + un.Username + "\n"
	}
	userString += "<b>About:</b> <code>" + uf.About + "</code>\n"
	if un.Usernames != nil {
		userString += "<b>Res. Usernames:</b> [<b>" + func() string {
			var s string
			for _, v := range un.Usernames {
				s += "@" + v.Username + " "
			}
			return s
		}() + "</b>]\n"
	}
	userString += "<b>User Link:</b> <a href=\"tg://user?id=" + strconv.FormatInt(un.ID, 10) + "\">userLink</a>\n\n---- User ID ----\n\n<b>ID:</b> " + strconv.FormatInt(un.ID, 10) + "\n"

	if uf.ProfilePhoto != nil {
		p := uf.ProfilePhoto.(*telegram.PhotoObj)
		if uf.PersonalPhoto != nil {
			p = uf.PersonalPhoto.(*telegram.PhotoObj)
		}
		inp := &telegram.InputMediaPhoto{
			ID: &telegram.InputPhotoObj{
				ID:            p.ID,
				AccessHash:    p.AccessHash,
				FileReference: p.FileReference,
			},
			Spoiler: true,
		}
		_, err := m.ReplyMedia(inp, telegram.MediaOptions{Caption: userString})
		if err != nil {
			m.Reply("Error: " + err.Error())
		}
	} else {
		m.Reply(userString)
	}
	return nil
}

func JsonHandle(m *telegram.NewMessage) error {
	if !m.IsReply() {
		if strings.Contains(m.Args(), "-s") {
			b, _ := json.MarshalIndent(m.Sender, "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else if strings.Contains(m.Args(), "-m") {
			b, _ := json.MarshalIndent(m.Media(), "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else if strings.Contains(m.Args(), "-c") {
			b, _ := json.MarshalIndent(m.Channel, "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else {
			m.Reply(`<pre lang="json">` + m.Marshal() + `</pre>`)
		}
	} else {
		r, _ := m.GetReplyMessage()
		if strings.Contains(m.Args(), "-s") {
			b, _ := json.MarshalIndent(r.Sender, "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else if strings.Contains(m.Args(), "-m") {
			b, _ := json.MarshalIndent(r.Media(), "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else if strings.Contains(m.Args(), "-c") {
			b, _ := json.MarshalIndent(r.Channel, "", "  ")
			m.Reply(`<pre lang="json">` + string(b) + `</pre>`)
		} else {
			m.Reply(`<pre lang="json">` + r.Marshal() + `</pre>`)
		}
	}

	return nil
}

func PingHandle(m *telegram.NewMessage) error {
	t := time.Now()
	x, _ := m.Reply("Pong!")
	x.Edit(fmt.Sprintf("Pong! <code>%s</code>", time.Since(t).String()))
	return nil
}
