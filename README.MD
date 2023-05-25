# ImageResizer

- Build Docker Image

```bash
docker build . -t resizer
```

- Run Docker Container
```bash
docker run --name resizer-dev -p 8080:8080 resizer n
```
where **n** number of concurrent requests

- Request
```bash
curl '127.0.0.1:8080/?width=200&height=200&url=http%3A%2F%2Fimages.pexels.com%2Fphotos%2F268533%2Fpexels-photo-268533.jpeg'
```