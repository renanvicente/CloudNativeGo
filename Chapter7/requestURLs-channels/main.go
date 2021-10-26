package main

type Resource struct {
	url			string
	polling		bool
	lastPooled	int64
}

func Pooler(in, out chan *Resource, numPollers int)  {
	for i := 0; i < numPollers; i++ {
		go func() {
			for r := range in {
				// Poll the URL

				// Send the processed Resource to out
				out <- r
			}
		}()
	}
}