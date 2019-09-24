package services

import (
	"bufio"
	"github.com/mgranderath/SPaaS/common"
	"github.com/mgranderath/SPaaS/server/model"
	"io"
	"log"
)

func (a *AppService) GetLogsStream(name string, messages model.StatusChannel) {
	response, err := a.DockerClient.ContainerLogs(common.SpaasName(name))
	if err != nil {
		messages.SendError(err)
		close(messages)
		return
	}
	defer response.Close()
	rd := bufio.NewReader(response)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("read file line error: %v", err)
			messages.SendError(err)
			close(messages)
			return
		}
		messages.SendInfo(line)
	}
	close(messages)
	return
}
