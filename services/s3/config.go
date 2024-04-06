package s3

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
}

func GetClient(accessKey, secretKey, region string) *S3Client {
	opts := s3.Options{
		Region:      *aws.String(region),
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "IQoJb3JpZ2luX2VjEJD//////////wEaCXVzLXdlc3QtMiJHMEUCIFlBdbKFNyK1oq2sf3iOuSxe2Nt3Jnq/Op6himIwENWlAiEA0LbCLILt1zRX3Qj2ePhrfnFah5BqmNsfg5vlL0ERLH4qoAMIuf//////////ARAAGgw2NTQ2NTQ0NTEzOTAiDO5zIq55pmnA895/5yr0AnyHhVhEXk+U9BGXNJk30ACSdRiykJbTiF9Ept8XGQrMr+7hbGQOXrz91aphW0DkKNWsg0qlCCzZXnfBzSVey5wGyedK2cDBF1K40y36cp47X0LGT6bCUcEhfs9JAc/utUA2ZZ0nQLQBf+tOEntMf2SKro0cI7gmYV9OGh/KdbM8Le8RU0u+z3owXgWGiV5P78LZL9FFGL7il8qBhokgt5WYIGa22DJb9EIPHKQOEk+LTB4NAHiWun1fNS9rEv/UCO4i2gV2Vam0zScdWkzYGY4S3szPpKDIrd0u4EzgDw/bBobbtTJ/vmLFfG8xtLE2h62Cc2uo3IkXnG8JpPrnporj5nHyPd7hePzP6r2MgzAellOI0eNQ5zngMC9xUVXd4rPrltV6Ts7Y6HQoBJiIHh75+tenBBE5hQQXHXSaU5TcPL7zlEpDPxz1tQIMp/o+2bwU1haV40B3iEJeE4ZqoYTf0odBX2DtxkcaUWTwsQcUq8SP4zDi9cOwBjqmAYOFxfqEOaxMZffMMS6OcPemeY4GbcVdBnHMXMyz2dBDCGmHun01hrZJnrizrai/huUYdsInwAHtUnB5D/JYAmfEt8rNQr2JqG1woJ9CwDjWriie5ArB621q/HCiCYyWTlzrktmiEfla4F/3Tkpwke5+Y6VIrLv2pxn2HcaRB47NSVIptkINNd1uS/ewM5vVAyrpgIv+GlOSp1ygvzeGEnYeXEVmZ/w=")),
	}

	// Create an Amazon S3 service client
	s3Client := s3.New(opts)

	client := &S3Client{
		Client: s3Client,
	}

	return client
}
