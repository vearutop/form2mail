package usecase

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/bool64/ctxd"
	"github.com/bool64/stats"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"github.com/vearutop/form2mail/internal/infra/email"
	"github.com/vearutop/form2mail/internal/infra/recaptcha"
)

type receiveDeps interface {
	CtxdLogger() ctxd.Logger
	StatsTracker() stats.Tracker
	RecaptchaChecker() recaptcha.Checker
	EmailSender() email.Sender
}

type Form struct {
	GRecaptchaResponse string                  `formData:"g-recaptcha-response"`
	SuccessURL         string                  `formData:"success_url"`
	FailURL            string                  `formData:"fail_url"`
	Values             []string                `formData:"values"`
	Files              []*multipart.FileHeader `formData:"files"`

	values map[string]string
	files  map[string]*multipart.FileHeader
	r      *http.Request
}

func (f *Form) LoadFromHTTPRequest(r *http.Request) error {
	isMultipart := r.Header.Get("Content-Type") == "multipart/form-data"

	if isMultipart {
		if err := r.ParseMultipartForm(30e6); err != nil {
			return err
		}

		for k, v := range r.MultipartForm.File {
			f.files[k] = v[0]
		}
	} else {
		r.ParseForm()
	}

	f.values = make(map[string]string, len(r.Form))
	f.files = make(map[string]*multipart.FileHeader)

	for k, v := range r.Form {
		f.values[k] = v[0]
	}

	f.GRecaptchaResponse = f.values["g-recaptcha-response"]
	delete(f.values, "g-recaptcha-response")

	f.SuccessURL = f.values["success_url"]
	delete(f.values, "success_url")

	f.FailURL = f.values["fail_url"]
	delete(f.values, "fail_url")

	f.r = r

	return nil
}

// Receive creates use case interactor.
func Receive(deps receiveDeps) usecase.Interactor {
	type receiveOutput struct {
		usecase.OutputWithEmbeddedWriter
	}

	u := usecase.NewInteractor(func(ctx context.Context, in Form, out *receiveOutput) (err error) {
		deps.CtxdLogger().Info(ctx, "", "in", in)

		w := out.Writer.(http.ResponseWriter)

		defer func() {
			if err != nil {
				if in.FailURL != "" {
					http.Redirect(w, in.r, in.FailURL, http.StatusSeeOther)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					_, _ = w.Write([]byte("ERROR: " + err.Error()))
				}
			}
		}()

		err = deps.RecaptchaChecker().CheckToken(ctx, in.GRecaptchaResponse)
		if err != nil {
			return err
		}

		err = deps.EmailSender().Send(ctx, in.values, in.files)
		if err != nil {
			return err
		}

		if in.SuccessURL != "" {
			http.Redirect(w, in.r, in.SuccessURL, http.StatusSeeOther)
			return nil
		}

		_, _ = w.Write([]byte("OK"))

		return nil
	})

	u.SetDescription("This endpoint receives HTTP form and forwards data to email.")
	u.SetTags("Forward")
	u.SetExpectedErrors(status.Unknown)

	return u
}
