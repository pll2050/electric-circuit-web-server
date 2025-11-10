using Microsoft.AspNetCore.Mvc;
using ElectricCircuitWeb.API.Services;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 인증 컨트롤러
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class AuthController : ControllerBase
{
    private readonly IAuthService _authService;
    private readonly ILogger<AuthController> _logger;

    public AuthController(IAuthService authService, ILogger<AuthController> logger)
    {
        _authService = authService;
        _logger = logger;
    }

    /// <summary>
    /// Firebase ID 토큰을 검증하고 사용자 정보를 반환합니다.
    /// </summary>
    [HttpPost("verify")]
    public async Task<IActionResult> VerifyToken([FromBody] VerifyTokenRequest request)
    {
        try
        {
            var firebaseUid = await _authService.VerifyIdTokenAsync(request.IdToken);
            var user = await _authService.GetUserByFirebaseUidAsync(firebaseUid);

            if (user == null)
            {
                return NotFound(new { message = "User not found" });
            }

            return Ok(new
            {
                user = new
                {
                    user.Id,
                    user.FirebaseUid,
                    user.Email,
                    user.DisplayName
                }
            });
        }
        catch (UnauthorizedAccessException)
        {
            return Unauthorized(new { message = "Invalid token" });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error verifying token");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }

    /// <summary>
    /// 새 사용자를 생성합니다.
    /// </summary>
    [HttpPost("signup")]
    public async Task<IActionResult> SignUp([FromBody] SignUpRequest request)
    {
        try
        {
            var firebaseUid = await _authService.VerifyIdTokenAsync(request.IdToken);

            // 이미 존재하는 사용자인지 확인
            var existingUser = await _authService.GetUserByFirebaseUidAsync(firebaseUid);
            if (existingUser != null)
            {
                return Conflict(new { message = "User already exists" });
            }

            var user = await _authService.CreateUserAsync(
                firebaseUid,
                request.Email,
                request.DisplayName
            );

            return Ok(new
            {
                user = new
                {
                    user.Id,
                    user.FirebaseUid,
                    user.Email,
                    user.DisplayName
                }
            });
        }
        catch (UnauthorizedAccessException)
        {
            return Unauthorized(new { message = "Invalid token" });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error signing up");
            return StatusCode(500, new { message = "Internal server error" });
        }
    }
}

public record VerifyTokenRequest(string IdToken);
public record SignUpRequest(string IdToken, string Email, string DisplayName);
