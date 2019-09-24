package model

import (
	"github.com/pingcap/tidb/util/hack"
	"time"
)

// MysqlQueryPiece 查询信息
type MysqlQueryPiece struct {
	BaseQueryPiece

	SessionID    *string `json:"cid"`
	ClientHost   *string `json:"-"`
	ClientPort   int     `json:"-"`

	VisitUser    *string `json:"user"`
	VisitDB      *string `json:"db"`
	QuerySQL     *string `json:"sql"`
	CostTimeInMS int64   `json:"cms"`
}

type PooledMysqlQueryPiece struct {
	MysqlQueryPiece
	recoverPool *mysqlQueryPiecePool
	sliceBufferPool *sliceBufferPool
}

func NewPooledMysqlQueryPiece(
	sessionID, clientIP, visitUser, visitDB, clientHost, serverIP *string,
	clientPort, serverPort int, throwPacketRate float64, stmtBeginTime int64) (
	mqp *PooledMysqlQueryPiece) {
	mqp = mqpp.Dequeue()

	nowInMS := time.Now().UnixNano() / millSecondUnit
	mqp.SessionID = sessionID
	mqp.ClientHost = clientIP
	mqp.ClientPort = clientPort
	mqp.ClientHost = clientHost
	mqp.ServerIP = serverIP
	mqp.ServerPort = serverPort
	mqp.VisitUser = visitUser
	mqp.VisitDB = visitDB
	mqp.SyncSend = false
	mqp.ThrowPacketRate = throwPacketRate
	mqp.BeginTime = stmtBeginTime
	mqp.CostTimeInMS = nowInMS - stmtBeginTime
	mqp.recoverPool = mqpp
	mqp.sliceBufferPool = localSliceBufferPool

	return
}

func (mqp *MysqlQueryPiece) String() (*string) {
	content := mqp.Bytes()
	contentStr := hack.String(content)
	return &contentStr
}

func (mqp *MysqlQueryPiece) Bytes() (content []byte) {
	// content, err := json.Marshal(mqp)
	if len(mqp.jsonContent) > 0 {
		return mqp.jsonContent
	}

	mqp.jsonContent = marsharQueryPiece(mqp)
	return mqp.jsonContent
}

func (mqp *MysqlQueryPiece) GetSQL() (str *string) {
	return mqp.QuerySQL
}

func (pmqp *PooledMysqlQueryPiece) Recovery() {
	pmqp.recoverPool.Enqueue(pmqp)
	pmqp.sliceBufferPool.Enqueue(pmqp.jsonContent[:0])
	pmqp.jsonContent = nil
}
