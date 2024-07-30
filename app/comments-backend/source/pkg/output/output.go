package output

import (
	"fmt"
	"io"
)

type msgType int

const (
	msgLog msgType = iota
	msgErr
)

type msg struct {
	t msgType
	m any
}

type Output struct {
	chOut chan<- msg
}

// Создаёт канал пересылки сообщений и логов. Все сообщения
// будут потокобезопасно записаны в log, а ошибки - в err.
func Make(log, err io.WriteCloser) Output {
	ch := make(chan msg)
	go loopOutput(ch, log, err)
	return Output{chOut: ch}
}

// Закрывает канал пересылки ошибок и логов. Текущие
// (ожидающие в других горутинах) и дальшейшие вызовы других
// методов не смогут отправить сообщение и будут возвращать
// false.
func (o Output) Close() {
	close(o.chOut)
}

// Отправляет сообщение в канал log. Если сообщение было
// доставлено, то вернет true, иначе - false
func (o Output) Log(message any) (ok bool) {
	return o.send(msgLog, message)
}

// Отправляет сообщение в канал err. Если сообщение было
// доставлено, то вернет true, иначе - false
func (o Output) Err(message any) (ok bool) {
	return o.send(msgErr, message)
}

func (o Output) send(t msgType, m any) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	o.chOut <- msg{t: t, m: m}
	return true
}

func loopOutput(chIn <-chan msg, log, err io.WriteCloser) {
	defer func() {
		log.Close()
		err.Close()
	}()

	for m := range chIn {
		switch m.t {
		case msgLog:
			fmt.Fprintln(log, m.m)
		case msgErr:
			fmt.Fprintln(err, m.m)
		}
	}
}
