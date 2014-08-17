package main

import (
	"time"
)

func MessageDispatcher(input chan Message) {
	for Running {
		select {
		case <-time.After(1e9):
			continue

		case msg := <-input:
			switch msg.(type) {
			case DirectoryRemoveMessage:

			case DirectoryCreateMessage:

			case DirectoryDiffMessage:

			case ShareDiscoverMessage:

			case ShareACKMessage:

			case ShareLeaveMessage:

			case ShareLastModMessage:

			case FileRemoveMessage:

			case FileUpdatedMessage:

			case FileCreatedMessage:

			case FileHashMessage:

			default:
				//Not implemented!!
			}
		}

	}
}
