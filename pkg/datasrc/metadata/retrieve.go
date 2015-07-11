package metadata

func GetMetaData() Digest {
	mch := make(chan Digest)
	defer close(mch)

	go func(mch chan Digest) {
		for _, p := range SupportedProviders {
			go func(c chan Digest) {
				m, err := p.MetaData("latest")
				if err != nil {
					return
				}

				c <- m.Digest()
			}(mch)
		}
	}(mch)

	return <-mch
}
