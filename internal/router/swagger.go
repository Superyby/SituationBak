package router

import (
	_ "embed"

	"SituationBak/docs"
	"github.com/gofiber/fiber/v3"
)

// Swagger UI HTML 模板
const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Orbital Tracker API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin: 0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: "/swagger/doc.json",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                persistAuthorization: true
            });
            window.ui = ui;
        };
    </script>
</body>
</html>`

// setupSwagger 配置 Swagger 路由
func setupSwagger(app *fiber.App) {
	// Swagger UI 页面
	app.Get("/swagger", func(c fiber.Ctx) error {
		return c.Redirect().To("/swagger/index.html")
	})

	app.Get("/swagger/index.html", func(c fiber.Ctx) error {
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(swaggerUIHTML)
	})

	// Swagger JSON 文档
	app.Get("/swagger/doc.json", func(c fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		return c.SendString(docs.SwaggerInfo.ReadDoc())
	})
}
