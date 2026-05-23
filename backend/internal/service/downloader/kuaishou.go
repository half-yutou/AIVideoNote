package downloader

type KuaiShouDownloader struct {
	GenericYtDlpDownloader
}

func NewKuaiShouDownloader() *KuaiShouDownloader {
	return &KuaiShouDownloader{
		GenericYtDlpDownloader: GenericYtDlpDownloader{platform: "kuaishou"},
	}
}
