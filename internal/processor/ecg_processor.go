package processor

import (
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
