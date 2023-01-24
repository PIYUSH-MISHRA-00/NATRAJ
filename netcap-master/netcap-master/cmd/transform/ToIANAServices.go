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

package transform

import (
	netmaltego "github.com/dreadl0ck/netcap/maltego"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/dreadl0ck/maltego"
	"github.com/dreadl0ck/netcap/resolvers"
	"github.com/dreadl0ck/netcap/types"
)

func toIANAServices() {
	resolverLog := zap.New(zapcore.NewNopCore())
	defer func() {
		err := resolverLog.Sync()
		if err != nil {
			log.Println(err)
		}
	}()

	resolvers.SetLogger(resolverLog)

	stdOut := os.Stdout
	os.Stdout = os.Stderr
	resolvers.InitServiceDB()
	os.Stdout = stdOut

	netmaltego.ConnectionTransform(
		netmaltego.CountOutgoingConnBytesFiltered,
		func(lt maltego.LocalTransform, trx *maltego.Transform, conn *types.Connection, min, max uint64, path string, mac string, ipaddr string, top12 *[]int) {
			i, err := strconv.Atoi(conn.DstPort)
			if err != nil {
				return
			}

			service := resolvers.LookupServiceByPort(i, strings.ToLower(conn.TransportProto))
			if service != "" {
				ent := addEntityWithPath(trx, "netcap.Service", service, path)
				ent.SetLinkLabel(humanize.Bytes(uint64(conn.TotalSize)))
				ent.SetLinkThickness(maltego.GetThickness(uint64(conn.TotalSize), min, max))
				ent.AddProperty("ip", "IP", maltego.Loose, conn.DstIP)
				ent.AddProperty("port", "Port", maltego.Loose, conn.DstPort)
			}
		},
	)
}
