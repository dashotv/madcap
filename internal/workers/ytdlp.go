package workers

// YtdlpJob runs a yt-dlp and wraps the output in a job
func (c *Client) YtdlpJob(url string) Job {
	return c.CommandJob("yt-dlp", url)
}

// YtdlpWithOptionsJob runs a yt-dlp with additional command line options and wraps the output in a job
func (c *Client) YtdlpWithOptionsJob(url string) Job {
	// TODO: handle custom options
	return c.CommandJob("yt-dlp", url)
}
