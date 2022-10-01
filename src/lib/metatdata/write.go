package metadata

func WriteLyrics(filepath string, contents string) error {
	return runMetadataPy("lyrics", filepath, contents)
}

func WriteArtist(filepath string, contents string) error {
	return runMetadataPy("artist", filepath, contents)
}

func WriteTitle(filepath string, contents string) error {
	return runMetadataPy("title", filepath, contents)
}

func WriteClean(filepath string) error {
	return runMetadataPy("clean", filepath)
}
