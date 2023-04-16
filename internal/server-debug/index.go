package serverdebug

import (
	"html/template"

	"github.com/labstack/echo/v4"

	"github.com/razielsd/chat-service/internal/logger"
)

type page struct {
	Path        string
	Description string
}

type indexPage struct {
	pages []page
}

func newIndexPage() *indexPage {
	return &indexPage{}
}

func (i *indexPage) addPage(path string, description string) {
	userPage := page{
		Path:        path,
		Description: description,
	}
	i.pages = append(i.pages, userPage)
}

func (i indexPage) handler(eCtx echo.Context) error {
	return template.Must(template.New("index").Parse(`<html>
	<title>Chat Service Debug</title>
<body>
	<h2>Chat Service Debug</h2>
	<ul>
		{{range .Pages }}
		<li><a href="{{ .Path }}">{{ .Description }}</a>
		{{end}}
	</ul>

	<h2>Log Level</h2>
	<form onSubmit="putLogLevel()">
		<select id="log-level-select">
		{{ $LogLevel := .LogLevel }}
		{{range $k,$v := .LogLevelList }}
		<option {{ if (eq $v $LogLevel) }} selected="selected" {{end}} >{{ $v }}</option>
		{{end}}
		</select>
		<input type="submit" value="Change"></input>
	</form>
	
	<script>
		function putLogLevel() {
			const req = new XMLHttpRequest();
			req.open('PUT', '/log/level', false);
			req.setRequestHeader('Content-type','application/x-www-form-urlencoded; charset=utf-8');
			req.onload = function() { window.location.reload(); };
			req.send('level='+document.getElementById('log-level-select').value);
		};
	</script>
</body>
</html>
`)).Execute(eCtx.Response(), struct {
		Pages        []page
		LogLevel     string
		LogLevelList []string
	}{
		Pages:        i.pages,
		LogLevelList: logger.LogLevelList,
		LogLevel:     logger.LogLevel.String(),
	})
}
