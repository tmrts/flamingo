package metadata

func GetMetaData() Digest {
	mch := make(chan Digest)
	defer close(mch)

	go func(mch chan Digest) {
		for t, p := range SupportedProviders {
			go func(c chan Digest) {
				c <- p.Metadata("latest").Digest()
			}(mch)
		}
	}(mch)

	return <-mch
}
