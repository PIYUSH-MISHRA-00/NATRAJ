/*
 * NETCAP - Traffic Analysis Framework
 * Copyright (c) 2017-2020 Philipp Mieden <dreadl0ck [at] protonmail [dot] ch>
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

package types

import (
	"github.com/dreadl0ck/netcap/encoder"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const fieldSkipCount = "SkipCount"

var fieldsEthernetCTP = []string{
	fieldTimestamp,
	fieldSkipCount, // int32
}

// CSVHeader returns the CSV header for the audit record.
func (i *EthernetCTP) CSVHeader() []string {
	return filter(fieldsEthernetCTP)
}

// CSVRecord returns the CSV record for the audit record.
func (i *EthernetCTP) CSVRecord() []string {
	return filter([]string{
		formatTimestamp(i.Timestamp),
		formatInt32(i.SkipCount),
	})
}

// Time returns the timestamp associated with the audit record.
func (i *EthernetCTP) Time() int64 {
	return i.Timestamp
}

// JSON returns the JSON representation of the audit record.
func (i *EthernetCTP) JSON() (string, error) {
	// convert unix timestamp from nano to millisecond precision for elastic
	i.Timestamp /= int64(time.Millisecond)

	return jsonMarshaler.MarshalToString(i)
}

var ethernetCTPMetric = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: strings.ToLower(Type_NC_EthernetCTP.String()),
		Help: Type_NC_EthernetCTP.String() + " audit records",
	},
	fieldsEthernetCTP[1:],
)

// Inc increments the metrics for the audit record.
func (i *EthernetCTP) Inc() {
	ethernetCTPMetric.WithLabelValues(i.CSVRecord()[1:]...).Inc()
}

// SetPacketContext sets the associated packet context for the audit record.
func (a *EthernetCTP) SetPacketContext(*PacketContext) {}

// Src TODO.
// Src returns the source address of the audit record.
func (i *EthernetCTP) Src() string {
	return ""
}

// Dst returns the destination address of the audit record.
func (i *EthernetCTP) Dst() string {
	return ""
}

var ethernetCTPEncoder = encoder.NewValueEncoder()

// Encode will encode categorical values and normalize according to configuration
func (i *EthernetCTP) Encode() []string {
	return filter([]string{
		ethernetCTPEncoder.Int64(fieldTimestamp, i.Timestamp),
		ethernetCTPEncoder.Int32(fieldSkipCount, i.SkipCount),
	})
}

// Analyze will invoke the configured analyzer for the audit record and return a score.
func (i *EthernetCTP) Analyze() {}

// NetcapType returns the type of the current audit record
func (a *EthernetCTP) NetcapType() Type {
	return Type_NC_EthernetCTP
}
