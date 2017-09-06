package chronograph

func ToEvent(obj interface{}) Event {
	return Event(nil)
}

func ToTimeline(obj interface{}) Timeline {
	return nil
}

func makeTimelineEventsPath(timeline Timeline) [][]byte {
	return [][]byte{boltTimelineBucketName, []byte(timeline.ID())}
}

func makeTimelineEventsPathWithKey(timeline Timeline) ([][]byte, []byte) {
	return makeTimelineEventsPath(timeline), boltTimelineEventsBucketName
}

func trimBytes(a []byte) (res []byte) {
	var skip = true
	for _, v := range a {
		if skip && v == 0 {
			continue
		} else {
			skip = false
		}
		res = append(res, v)
	}
	return
}
