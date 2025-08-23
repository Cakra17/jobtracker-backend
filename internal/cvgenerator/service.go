package cvgenerator

import (
	"context"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type PDFGenerator struct{
  renderer *TemplateRenderer
}

func NewPdfGenerator(renderer *TemplateRenderer) PDFGenerator {
  return PDFGenerator{renderer: renderer} 
}

func (p *PDFGenerator) GeneratePDF(ctx context.Context, resume *Resume) ([]byte, error) {
  html, err := p.renderer.RenderHTMLWithCSS(resume)
  if err != nil {
    return nil, err
  }

  chromeCtx, cancel := chromedp.NewContext(ctx)
  defer cancel()

  var pdfByte []byte
  err = chromedp.Run(chromeCtx,
    chromedp.Navigate("about:blank"),
    chromedp.ActionFunc(func(ctx context.Context) error {
      frameTree, err := page.GetFrameTree().Do(ctx)
      if err != nil {
        return err
      }
      return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
    }),
    chromedp.WaitVisible("body", chromedp.ByQuery),
    chromedp.ActionFunc(func(ctx context.Context) error {
      var err error
      pdfByte, _, err = page.PrintToPDF().
        WithPrintBackground(true).
        WithPaperWidth(8.5).
        WithPaperHeight(11).
        WithMarginTop(0.4).
        WithMarginBottom(0.4).
        WithMarginLeft(0.4).
        WithMarginRight(0.4).
        WithPreferCSSPageSize(false).
        Do(ctx) 
      return err
    }),
  )

  if err != nil {
    return nil, err
  }

  return pdfByte, nil
}

