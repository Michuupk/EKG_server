package processor

import (
	"context"
	"ekg_server/internal/domain"
	"time"
)

func UnpackPacket(packet domain.ECGPacket, baseTime time.Time) []domain.ECGDataPoint {
	points := make([]domain.ECGDataPoint, 0, len(packet.Packet))

	for i, val := range packet.Packet {
		point := domain.ECGDataPoint{
			Timestamp: baseTime.Add(time.Duration(i) * 4 * time.Millisecond),
			DeviceID:  packet.DeviceID,
			PatientID: packet.PatientID,
			Seq:       packet.Seq,
			Value:     val,
			IsAnomaly: val < 0 || val > 1023,
		}
		points = append(points, point)
	}

	return points
}

func DispatchPoints(ctx context.Context, points []domain.ECGDataPoint, storageChan chan<- domain.ECGDataPoint, alertChan chan<- domain.ECGDataPoint) error {
	for _, pt := range points {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case storageChan <- pt:
		}

		if pt.IsAnomaly {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case alertChan <- pt:
			}
		}
	}
	return nil
}

func StartProcessor(
	ctx context.Context,
	inChan <-chan domain.ECGPacket,
	storageChan chan<- domain.ECGDataPoint,
	alertChan chan<- domain.ECGDataPoint,
) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pkt, ok := <-inChan:
			if !ok {
				return nil
			}
			if err := pkt.Validate(); err != nil {
				continue
			}

			points := UnpackPacket(pkt, time.Now())

			if err := DispatchPoints(ctx, points, storageChan, alertChan); err != nil {
				return err
			}
		}
	}
}
