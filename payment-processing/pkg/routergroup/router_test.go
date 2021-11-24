package routergroup

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandle(t *testing.T) {
	router := New()

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	router.Handle("GET", "/", h)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Error("Test Handle failed")
	}
}

func TestHandler(t *testing.T) {
	router := New()

	h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	router.Handler("GET", "/", h)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Error("Test Handler failed")
	}
}

func TestHandlerFunc(t *testing.T) {
	router := New()

	h := func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}
	router.HandlerFunc("GET", "/", h)

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	if w.Code != http.StatusTeapot {
		t.Error("Test HandlerFunc failed")
	}
}

func TestMethod(t *testing.T) {
	router := New()

	router.DELETE("/delete", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.GET("/get", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.HEAD("/head", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.OPTIONS("/options", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.PATCH("/patch", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.POST("/post", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	router.PUT("/put", http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))

	samples := map[string]string{
		"DELETE":  "/delete",
		"GET":     "/get",
		"HEAD":    "/head",
		"OPTIONS": "/options",
		"PATCH":   "/patch",
		"POST":    "/post",
		"PUT":     "/put",
	}
	for method, path := range samples {
		r := httptest.NewRequest(method, path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		if w.Code != http.StatusTeapot {
			t.Errorf("Path %s not registered", path)
		}
	}
}

func TestGroup(t *testing.T) {
	router := New()
	foo := router.Group("/foo")
	bar := router.Group("/bar")
	baz := foo.Group("/baz")

	foo.HandlerFunc("GET", "", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	foo.HandlerFunc("GET", "/group", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	bar.HandlerFunc("GET", "/group", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	baz.HandlerFunc("GET", "/group", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	samples := []string{"/foo", "/foo/group", "/foo/baz/group", "/bar/group"}

	for _, path := range samples {
		r := httptest.NewRequest("GET", path, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, r)
		if w.Code != http.StatusTeapot {
			t.Errorf("Grouped path %s not registered", path)
		}
	}
}

func TestMiddleware(t *testing.T) {
	var use, group bool

	router := New().Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			use = true
			next.ServeHTTP(w, r)
		})
	})

	foo := router.Group("/foo", func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			group = true
			next.ServeHTTP(w, r)
		})
	})

	foo.HandlerFunc("GET", "/bar", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})

	r := httptest.NewRequest("GET", "/foo/bar", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	if !use {
		t.Error("Middleware registered by Use() under \"/\" not touched")
	}
	if !group {
		t.Error("Middleware registered by Group() under \"/foo\" not touched")
	}
}

func TestStatic(t *testing.T) {
	files := []string{"temp_1", "temp_2"}
	strs := []string{"test content", "static contents"}

	for i := range files {
		f, _ := os.Create(files[i])
		defer os.Remove(files[i])

		f.WriteString(strs[i])
		f.Sync()
		f.Close()
	}

	pwd, _ := os.Getwd()
	router := New()
	router.Static("/*filepath", pwd)

	for i := range files {
		r := httptest.NewRequest("GET", "/"+files[i], nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, r)

		body := w.Result().Body
		defer body.Close()

		file, _ := ioutil.ReadAll(body)
		if string(file) != strs[i] {
			t.Error("Test Static failed")
		}
	}
}

func TestFile(t *testing.T) {
	str := "test_content"

	f, _ := os.Create("temp_file")
	defer os.Remove("temp_file")

	f.WriteString(str)
	f.Sync()
	f.Close()

	router := New()
	router.File("/file", "temp_file")

	r := httptest.NewRequest("GET", "/file", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)

	body := w.Result().Body
	defer body.Close()

	file, _ := ioutil.ReadAll(body)
	if string(file) != str {
		t.Error("Test File failed")
	}
}
