using FirebaseAdmin.Auth;

namespace ElectricCircuitWeb.API.Middleware;

/// <summary>
/// Firebase 인증 미들웨어
/// </summary>
public class FirebaseAuthMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<FirebaseAuthMiddleware> _logger;

    public FirebaseAuthMiddleware(RequestDelegate next, ILogger<FirebaseAuthMiddleware> logger)
    {
        _next = next;
        _logger = logger;
    }

    public async Task InvokeAsync(HttpContext context)
    {
        // 헬스 체크 및 인증 관련 엔드포인트는 스킵
        if (context.Request.Path.StartsWithSegments("/api/health") ||
            context.Request.Path.StartsWithSegments("/api/auth"))
        {
            await _next(context);
            return;
        }

        var authHeader = context.Request.Headers["Authorization"].FirstOrDefault();
        if (authHeader?.StartsWith("Bearer ") == true)
        {
            var token = authHeader.Substring("Bearer ".Length).Trim();

            try
            {
                var decodedToken = await FirebaseAuth.DefaultInstance.VerifyIdTokenAsync(token);
                context.Items["FirebaseUid"] = decodedToken.Uid;
                context.Items["FirebaseEmail"] = decodedToken.Claims.GetValueOrDefault("email");
            }
            catch (Exception ex)
            {
                _logger.LogWarning(ex, "Failed to verify Firebase token");
                context.Response.StatusCode = 401;
                await context.Response.WriteAsJsonAsync(new { message = "Unauthorized" });
                return;
            }
        }

        await _next(context);
    }
}

/// <summary>
/// 미들웨어 확장 메서드
/// </summary>
public static class FirebaseAuthMiddlewareExtensions
{
    public static IApplicationBuilder UseFirebaseAuth(this IApplicationBuilder builder)
    {
        return builder.UseMiddleware<FirebaseAuthMiddleware>();
    }
}
