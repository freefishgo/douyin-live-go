package player

import (
	"github.com/gopxl/beep"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/speaker"
	"time"
)

type Stream struct {
	Ctrl      *beep.Ctrl
	Volume    *effects.Volume
	Speedy    *beep.Resampler
	EndStream beep.Streamer
	WaitEnd   chan struct{}
}

func (s *Stream) IsPlaying() bool {
	select {
	case <-s.WaitEnd:
		return true
	default:
		return false
	}
}

func (s *Stream) WaitPlayerEnd() {
	<-s.WaitEnd
}

var stream *Stream

var queue = &Queue{
	cacheLen: make(chan struct{}, 1),
}

func init() {
	sr := beep.SampleRate(44100)
	speaker.Init(sr, sr.N(time.Second/10))
	stream = newStream(queue)
	speaker.Play(stream.EndStream)
}

func Play(streamer beep.Streamer, playEnd func()) {
	//speaker.Lock()
	//defer speaker.Unlock()
	queue.cacheLen <- struct{}{}
	ctr := newStream(streamer)
	go func() {
		ctr.WaitPlayerEnd()
		<-queue.cacheLen
		playEnd()
	}()
	queue.Add(ctr.EndStream)
	stream.Ctrl.Paused = false
	//if stream == nil {
	//	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	//	stream = newStream(streamer, loopCount)
	//} else {
	//	stream.WaitPlayerEnd()
	//	if playEnd != nil {
	//		playEnd()
	//	}
	//	stream = newStream(streamer, loopCount)
	//}
	//
	//speaker.Play(stream.EndStream)
	//if stream.Ctrl.Paused {
	//	stream.Ctrl.Paused = false
	//}
}

func newStream(streamer beep.Streamer) *Stream {
	//ctrl := &beep.Ctrl{Streamer: beep.Loop(loopCount, streamer), Paused: false}
	ctrl := &beep.Ctrl{Streamer: streamer, Paused: false}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speedy := beep.ResampleRatio(4, 1, volume)
	s := &Stream{
		Ctrl:    ctrl,
		Volume:  volume,
		Speedy:  speedy,
		WaitEnd: make(chan struct{}),
	}
	s.EndStream = beep.Seq(speedy, beep.Callback(func() {
		defer func() {
			recover()
		}()
		close(s.WaitEnd)
	}))
	return s
}

type Queue struct {
	streamers []beep.Streamer
	cacheLen  chan struct{}
}

func (q *Queue) Add(streamers ...beep.Streamer) {
	q.streamers = append(q.streamers, streamers...)
}

func (q *Queue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}
	return len(samples), true
}

func (q *Queue) Err() error {
	return nil
}
