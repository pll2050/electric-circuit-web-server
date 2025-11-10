using Microsoft.AspNetCore.Mvc;
using ElectricCircuitWeb.API.Services;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 사용자 관리 컨트롤러 (관리자용)
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    private readonly IAuthService _authService;
    private readonly ILogger<UsersController> _logger;

    public UsersController(IAuthService authService, ILogger<UsersController> logger)
    {
        _authService = authService;
        _logger = logger;
    }

    /// <summary>
    /// 모든 사용자 목록을 조회합니다.
    /// 먼저 Firebase Authentication에서 시도하고, 실패하면 PostgreSQL에서 조회합니다.
    /// GET /api/users
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> GetAllUsers()
    {
        try
        {
            IEnumerable<Models.User> users;

            try
            {
                // 먼저 Firebase Authentication에서 시도
                _logger.LogInformation("Attempting to fetch users from Firebase Authentication");
                users = await _authService.GetAllFirebaseUsersAsync();

                if (!users.Any())
                {
                    _logger.LogWarning("No users returned from Firebase, falling back to PostgreSQL");
                    users = await _authService.GetAllUsersAsync();
                }
            }
            catch (Exception firebaseEx)
            {
                _logger.LogError(firebaseEx, "Firebase user listing failed, falling back to PostgreSQL");
                // Fallback to PostgreSQL
                users = await _authService.GetAllUsersAsync();
            }

            return Ok(new
            {
                success = true,
                users = users.Select(u => new
                {
                    uid = u.FirebaseUid,
                    email = u.Email,
                    displayName = u.DisplayName,
                    photoURL = u.PhotoUrl,
                    phoneNumber = u.PhoneNumber,
                    provider = u.Provider,
                    createdAt = u.CreatedAt,
                    updatedAt = u.UpdatedAt,
                    lastLoginAt = u.LastLoginAt
                })
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error fetching all users");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }
}
