package gocbcore

import (
	"fmt"
	"net"
	"strings"
)

type routeConfig struct {
	revId        int64
	uuid         string
	bktType      bucketType
	kvServerList []string
	capiEpList   []string
	mgmtEpList   []string
	n1qlEpList   []string
	ftsEpList    []string
	vbMap        *vbucketMap
	ketamaMap    *ketamaContinuum
}

func (config *routeConfig) IsValid() bool {
	if len(config.kvServerList) == 0 || len(config.mgmtEpList) == 0 {
		return false
	}
	switch config.bktType {
	case bktTypeCouchbase:
		return config.vbMap != nil && config.vbMap.IsValid()
	case bktTypeMemcached:
		return config.ketamaMap != nil && config.ketamaMap.IsValid()
	default:
		return false
	}
}

func buildRouteConfig(bk *cfgBucket, useSsl bool) *routeConfig {
	var kvServerList []string
	var capiEpList []string
	var mgmtEpList []string
	var n1qlEpList []string
	var ftsEpList []string
	var bktType bucketType

	switch bk.NodeLocator {
	case "ketama":
		bktType = bktTypeMemcached
	case "vbucket":
		bktType = bktTypeCouchbase
	default:
		logDebugf("Invalid nodeLocator %s", bk.NodeLocator)
		bktType = bktTypeInvalid
	}

	if bk.NodesExt != nil {
		for _, node := range bk.NodesExt {
			// Hostname blank means to use the same one as was connected to
			if node.Hostname == "" {
				node.Hostname = bk.SourceHostname
			}

			if !useSsl {
				if node.Services.Kv > 0 {
					kvServerList = append(kvServerList, fmt.Sprintf("%s:%d", node.Hostname, node.Services.Kv))
				}
				if node.Services.Capi > 0 {
					capiEpList = append(capiEpList, fmt.Sprintf("http://%s:%d/%s", node.Hostname, node.Services.Capi, bk.Name))
				}
				if node.Services.Mgmt > 0 {
					mgmtEpList = append(mgmtEpList, fmt.Sprintf("http://%s:%d", node.Hostname, node.Services.Mgmt))
				}
				if node.Services.N1ql > 0 {
					n1qlEpList = append(n1qlEpList, fmt.Sprintf("http://%s:%d", node.Hostname, node.Services.N1ql))
				}
				if node.Services.Fts > 0 {
					ftsEpList = append(ftsEpList, fmt.Sprintf("http://%s:%d", node.Hostname, node.Services.Fts))
				}
			} else {
				if node.Services.KvSsl > 0 {
					kvServerList = append(kvServerList, fmt.Sprintf("%s:%d", node.Hostname, node.Services.KvSsl))
				}
				if node.Services.CapiSsl > 0 {
					capiEpList = append(capiEpList, fmt.Sprintf("https://%s:%d/%s", node.Hostname, node.Services.CapiSsl, bk.Name))
				}
				if node.Services.MgmtSsl > 0 {
					mgmtEpList = append(mgmtEpList, fmt.Sprintf("https://%s:%d", node.Hostname, node.Services.MgmtSsl))
				}
				if node.Services.N1qlSsl > 0 {
					n1qlEpList = append(n1qlEpList, fmt.Sprintf("https://%s:%d", node.Hostname, node.Services.N1qlSsl))
				}
				if node.Services.FtsSsl > 0 {
					ftsEpList = append(ftsEpList, fmt.Sprintf("http://%s:%d", node.Hostname, node.Services.FtsSsl))
				}
			}
		}
	} else {
		if useSsl {
			logErrorf("Received config without nodesExt while SSL is enabled.  Generating invalid config.")
			return &routeConfig{}
		}

		if bktType == bktTypeCouchbase {
			kvServerList = bk.VBucketServerMap.ServerList
		}

		for _, node := range bk.Nodes {
			if node.CouchAPIBase != "" {
				// Slice off the UUID as Go's HTTP client cannot handle being passed URL-Encoded path values.
				capiEp := strings.SplitN(node.CouchAPIBase, "%2B", 2)[0]

				capiEpList = append(capiEpList, capiEp)
			}
			if node.Hostname != "" {
				mgmtEpList = append(mgmtEpList, fmt.Sprintf("http://%s", node.Hostname))
			}

			if bktType == bktTypeMemcached {
				// Get the data port. No VBucketServerMap.
				host, _, err := net.SplitHostPort(node.Hostname)
				if err != nil {
					logErrorf("Encountered invalid memcached host/port string. Ignoring node.")
					continue
				}

				curKvHost := fmt.Sprintf("%s:%d", host, node.Ports["direct"])
				kvServerList = append(kvServerList, curKvHost)
			}
		}
	}

	rc := &routeConfig{
		revId:        bk.Rev,
		uuid:         bk.UUID,
		kvServerList: kvServerList,
		capiEpList:   capiEpList,
		mgmtEpList:   mgmtEpList,
		n1qlEpList:   n1qlEpList,
		ftsEpList:    ftsEpList,
		bktType:      bktType,
	}

	if bktType == bktTypeCouchbase {
		vbMap := bk.VBucketServerMap.VBucketMap
		numReplicas := bk.VBucketServerMap.NumReplicas
		rc.vbMap = newVbucketMap(vbMap, numReplicas)
	} else if bktType == bktTypeMemcached {
		rc.ketamaMap = newKetamaContinuum(kvServerList)
	}

	return rc
}
