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

    // ==================== 토큰 검증 ====================

    /// <summary>
    /// Firebase ID 토큰을 검증하고 사용자 정보를 반환합니다.
    /// POST /api/auth/verify
    /// </summary>
    [HttpPost("verify")]
    public async Task<IActionResult> VerifyToken([FromBody] VerifyTokenRequest request)
    {
        try
        {
            var firebaseUid = await _authService.VerifyIdTokenAsync(request.Token);
            var firebaseUser = await _authService.GetFirebaseUserAsync(firebaseUid);

            if (firebaseUser == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "User not found"
                });
            }

            return Ok(new
            {
                success = true,
                message = "Token verified successfully",
                user = new
                {
                    uid = firebaseUser.Uid,
                    email = firebaseUser.Email,
                    emailVerified = firebaseUser.EmailVerified,
                    displayName = firebaseUser.DisplayName
                }
            });
        }
        catch (UnauthorizedAccessException)
        {
            return Unauthorized(new
            {
                success = false,
                error = "Invalid token"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error verifying token");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 사용자 생성 ====================

    /// <summary>
    /// 서버에서 직접 Firebase 사용자를 생성합니다 (관리자용).
    /// POST /api/auth/create-user
    /// </summary>
    [HttpPost("create-user")]
    public async Task<IActionResult> CreateUser([FromBody] CreateUserRequest request)
    {
        try
        {
            if (string.IsNullOrEmpty(request.Email))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Email is required"
                });
            }

            if (string.IsNullOrEmpty(request.Password))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Password is required"
                });
            }

            // Firebase에 사용자 생성
            var firebaseUser = await _authService.CreateFirebaseUserAsync(
                request.Email,
                request.Password,
                request.DisplayName,
                request.PhotoUrl
            );

            // PostgreSQL DB에 사용자 정보 저장
            var user = await _authService.CreateUserAsync(
                firebaseUser.Uid,
                firebaseUser.Email,
                firebaseUser.DisplayName,
                firebaseUser.PhotoUrl,
                "password" // provider: 이메일/비밀번호 방식
            );

            return Ok(new
            {
                success = true,
                message = "User created successfully",
                user = new
                {
                    id = user.Id,
                    uid = firebaseUser.Uid,
                    email = firebaseUser.Email,
                    displayName = firebaseUser.DisplayName
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating user");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    /// <summary>
    /// 새 사용자를 가입시킵니다 (토큰 검증 후 DB에 저장).
    /// POST /api/auth/signup
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
                // 기존 사용자의 경우 LastLoginAt만 업데이트 (UpdatedAt은 변경하지 않음)
                var updatedUser = await _authService.UpdateLastLoginAtAsync(existingUser.Id);

                return Ok(new
                {
                    success = true,
                    message = "User login successful",
                    user = new
                    {
                        updatedUser.Id,
                        updatedUser.FirebaseUid,
                        updatedUser.Email,
                        updatedUser.DisplayName
                    }
                });
            }

            // Firebase에서 사용자 정보 가져오기
            var firebaseUser = await _authService.GetFirebaseUserAsync(firebaseUid);

            var user = await _authService.CreateUserAsync(
                firebaseUid,
                request.Email,
                request.DisplayName,
                firebaseUser?.PhotoUrl,
                request.Provider ?? "email" // 클라이언트가 전달한 provider 사용 (기본값: email)
            );

            return Ok(new
            {
                success = true,
                message = "User registered successfully",
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
            return Unauthorized(new
            {
                success = false,
                error = "Invalid token"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error signing up");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 사용자 조회 ====================

    /// <summary>
    /// 사용자 정보를 조회합니다.
    /// GET /api/auth/get-user?uid={uid}
    /// </summary>
    [HttpGet("get-user")]
    public async Task<IActionResult> GetUser([FromQuery] string uid)
    {
        try
        {
            if (string.IsNullOrEmpty(uid))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "User UID is required"
                });
            }

            var firebaseUser = await _authService.GetFirebaseUserAsync(uid);
            if (firebaseUser == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "User not found"
                });
            }

            return Ok(new
            {
                success = true,
                message = "User found",
                user = new
                {
                    uid = firebaseUser.Uid,
                    email = firebaseUser.Email,
                    displayName = firebaseUser.DisplayName,
                    photoURL = firebaseUser.PhotoUrl
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting user: {Uid}", uid);
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 사용자 수정 ====================

    /// <summary>
    /// 사용자 정보를 수정합니다.
    /// PUT /api/auth/update-user
    /// </summary>
    [HttpPut("update-user")]
    public async Task<IActionResult> UpdateUser([FromBody] UpdateUserRequest request)
    {
        try
        {
            if (string.IsNullOrEmpty(request.Uid))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "User UID is required"
                });
            }

            var firebaseUser = await _authService.UpdateFirebaseUserAsync(
                request.Uid,
                request.Email,
                request.DisplayName,
                request.PhotoUrl,
                request.Password,
                request.PhoneNumber
            );

            return Ok(new
            {
                success = true,
                message = "User updated successfully",
                user = new
                {
                    uid = firebaseUser.Uid,
                    displayName = firebaseUser.DisplayName,
                    photoURL = firebaseUser.PhotoUrl,
                    phoneNumber = firebaseUser.PhoneNumber
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating user: {Uid}", request.Uid);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    /// <summary>
    /// 사용자 프로필 정보를 수정합니다.
    /// PUT /api/auth/update-profile
    /// </summary>
    [HttpPut("update-profile")]
    public async Task<IActionResult> UpdateProfile([FromBody] UpdateProfileRequest request)
    {
        try
        {
            if (string.IsNullOrEmpty(request.IdToken))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "ID token is required"
                });
            }

            // ID 토큰 검증
            var firebaseUid = await _authService.VerifyIdTokenAsync(request.IdToken);

            // Firebase 사용자 정보 업데이트
            var firebaseUser = await _authService.UpdateFirebaseUserAsync(
                firebaseUid,
                null, // Email은 변경하지 않음
                request.DisplayName,
                request.PhotoUrl,
                null, // Password는 변경하지 않음
                request.PhoneNumber
            );

            // PostgreSQL DB의 사용자 정보도 업데이트
            var dbUser = await _authService.GetUserByFirebaseUidAsync(firebaseUid);
            if (dbUser != null)
            {
                dbUser.DisplayName = request.DisplayName ?? dbUser.DisplayName;
                dbUser.PhotoUrl = request.PhotoUrl ?? dbUser.PhotoUrl;
                dbUser.PhoneNumber = request.PhoneNumber ?? dbUser.PhoneNumber;
                await _authService.UpdateUserAsync(dbUser);
            }

            return Ok(new
            {
                success = true,
                message = "Profile updated successfully",
                user = new
                {
                    uid = firebaseUser.Uid,
                    displayName = firebaseUser.DisplayName,
                    photoURL = firebaseUser.PhotoUrl,
                    phoneNumber = firebaseUser.PhoneNumber
                }
            });
        }
        catch (UnauthorizedAccessException)
        {
            return Unauthorized(new
            {
                success = false,
                error = "Invalid token"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating profile");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 사용자 삭제 ====================

    /// <summary>
    /// 사용자를 삭제합니다.
    /// DELETE /api/auth/delete-user?uid={uid}
    /// </summary>
    [HttpDelete("delete-user")]
    public async Task<IActionResult> DeleteUser([FromQuery] string uid)
    {
        try
        {
            if (string.IsNullOrEmpty(uid))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "User UID is required"
                });
            }

            await _authService.DeleteFirebaseUserAsync(uid);
            await _authService.DeleteUserAsync(uid);

            return Ok(new
            {
                success = true,
                message = "User deleted successfully"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting user: {Uid}", uid);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 커스텀 클레임 ====================

    /// <summary>
    /// 사용자에게 커스텀 클레임을 설정합니다.
    /// POST /api/auth/set-custom-claims
    /// </summary>
    [HttpPost("set-custom-claims")]
    public async Task<IActionResult> SetCustomClaims([FromBody] SetCustomClaimsRequest request)
    {
        try
        {
            if (string.IsNullOrEmpty(request.Uid))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "User UID is required"
                });
            }

            if (request.CustomClaims == null)
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Custom claims are required"
                });
            }

            await _authService.SetCustomClaimsAsync(request.Uid, request.CustomClaims);

            return Ok(new
            {
                success = true,
                message = "Custom claims set successfully"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error setting custom claims: {Uid}", request.Uid);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }
}

// ==================== Request/Response Models ====================

public record VerifyTokenRequest(string Token);
public record SignUpRequest(string IdToken, string Email, string DisplayName, string? Provider = null);
public record CreateUserRequest(string Email, string Password, string? DisplayName = null, string? PhotoUrl = null);
public record UpdateUserRequest(string Uid, string? Email = null, string? DisplayName = null, string? PhotoUrl = null, string? Password = null, string? PhoneNumber = null);
public record UpdateProfileRequest(string IdToken, string? DisplayName = null, string? PhotoUrl = null, string? PhoneNumber = null);
public record SetCustomClaimsRequest(string Uid, Dictionary<string, object> CustomClaims);
