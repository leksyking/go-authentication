package wait

import "sync"

var Wg *sync.WaitGroup = new(sync.WaitGroup)
