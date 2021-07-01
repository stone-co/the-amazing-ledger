package ledger

import (
	"context"
	"fmt"

	"github.com/stone-co/the-amazing-ledger/app/domain/vos"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger"
	"google.golang.org/grpc/status"
)

type Path struct {
	Account string
	Debit   int64
	Credit  int64
}

type SyntheticReport struct {
	totalCredit int
	totalDebit  int
	paths       []*Path
}

func (a SyntheticReport) TotalCredit() int {
	return a.totalCredit
}

func (a SyntheticReport) TotalDebit() int {
	return a.totalDebit
}

func (a SyntheticReport) Paths() []*Path {
	return a.paths
}

func (c *Connection) GetSyntheticReport(ctx context.Context, accountName string, startTime int64, endTime int64) (*SyntheticReport, error) {

	accountPath, err := vos.NewAccountPath(accountName)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	reportRequest := &proto.GetSyntheticReportRequest{
		AccountPath: accountPath.Name(),
		StartTime:   startTime,
		EndTime:     endTime,
	}

	response, err := c.client.GetSyntheticReport(ctx, reportRequest)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			return nil, fmt.Errorf(e.Message())
		}

		return nil, fmt.Errorf("not able to parse error returned %v", err)
	}

	fmt.Printf("> responsee: %v", response)

	report := &SyntheticReport{
		totalCredit: int(response.TotalCredit),
		totalDebit:  int(response.TotalDebit),
		paths:       toPaths(response.Paths),
	}

	return report, nil
}

func toPaths(protoPaths []*proto.Path) []*Path {
	var paths []*Path = []*Path{}

	for _, v := range protoPaths {
		path := Path{
			Debit:   v.Debit,
			Credit:  v.Credit,
			Account: v.Account,
		}

		paths = append(paths, &path)
	}

	return paths
}
