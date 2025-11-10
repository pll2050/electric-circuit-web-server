using ElectricCircuitWeb.API.Models;
using ElectricCircuitWeb.API.Repositories;
using FirebaseAdmin.Auth;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 인증 서비스 구현
/// </summary>
public class AuthService : IAuthService
{
    private readonly IUserRepository _userRepository;
    private readonly ILogger<AuthService> _logger;

    public AuthService(IUserRepository userRepository, ILogger<AuthService> logger)
    {
        _userRepository = userRepository;
        _logger = logger;
    }

    // ==================== 토큰 검증 ====================

    public async Task<string> VerifyIdTokenAsync(string idToken)
    {
        try
        {
            var decodedToken = await FirebaseAuth.DefaultInstance.VerifyIdTokenAsync(idToken);
            return decodedToken.Uid;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to verify Firebase ID token");
            throw new UnauthorizedAccessException("Invalid token");
        }
    }

    // ==================== 사용자 조회 ====================

    public async Task<User?> GetUserByFirebaseUidAsync(string firebaseUid)
    {
        return await _userRepository.GetByFirebaseUidAsync(firebaseUid);
    }

    public async Task<UserRecord?> GetFirebaseUserAsync(string uid)
    {
        try
        {
            var user = await FirebaseAuth.DefaultInstance.GetUserAsync(uid);
            return user;
        }
        catch (FirebaseAuthException ex)
        {
            _logger.LogError(ex, "Failed to get Firebase user: {Uid}", uid);
            return null;
        }
    }

    // ==================== 사용자 생성 ====================

    public async Task<User> CreateUserAsync(string firebaseUid, string email, string displayName)
    {
        var user = new User
        {
            FirebaseUid = firebaseUid,
            Email = email,
            DisplayName = displayName
        };

        return await _userRepository.CreateAsync(user);
    }

    public async Task<UserRecord> CreateFirebaseUserAsync(string email, string password, string? displayName = null, string? photoUrl = null)
    {
        try
        {
            var args = new UserRecordArgs
            {
                Email = email,
                Password = password,
                DisplayName = displayName,
                PhotoUrl = photoUrl,
                EmailVerified = false
            };

            var user = await FirebaseAuth.DefaultInstance.CreateUserAsync(args);
            _logger.LogInformation("Created Firebase user: {Uid}, {Email}", user.Uid, user.Email);
            return user;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to create Firebase user");
            throw;
        }
    }

    // ==================== 사용자 수정 ====================

    public async Task<User> UpdateUserAsync(User user)
    {
        return await _userRepository.UpdateAsync(user);
    }

    public async Task<UserRecord> UpdateFirebaseUserAsync(string uid, string? email = null, string? displayName = null, string? photoUrl = null, string? password = null)
    {
        try
        {
            var args = new UserRecordArgs
            {
                Uid = uid
            };

            if (!string.IsNullOrEmpty(email))
                args.Email = email;
            if (!string.IsNullOrEmpty(displayName))
                args.DisplayName = displayName;
            if (!string.IsNullOrEmpty(photoUrl))
                args.PhotoUrl = photoUrl;
            if (!string.IsNullOrEmpty(password))
                args.Password = password;

            var user = await FirebaseAuth.DefaultInstance.UpdateUserAsync(args);
            _logger.LogInformation("Updated Firebase user: {Uid}", uid);
            return user;
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to update Firebase user: {Uid}", uid);
            throw;
        }
    }

    // ==================== 사용자 삭제 ====================

    public async Task<bool> DeleteUserAsync(string firebaseUid)
    {
        var user = await _userRepository.GetByFirebaseUidAsync(firebaseUid);
        if (user == null)
            return false;

        return await _userRepository.DeleteAsync(user.Id);
    }

    public async Task DeleteFirebaseUserAsync(string uid)
    {
        try
        {
            await FirebaseAuth.DefaultInstance.DeleteUserAsync(uid);
            _logger.LogInformation("Deleted Firebase user: {Uid}", uid);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to delete Firebase user: {Uid}", uid);
            throw;
        }
    }

    // ==================== 커스텀 클레임 ====================

    public async Task SetCustomClaimsAsync(string uid, Dictionary<string, object> customClaims)
    {
        try
        {
            await FirebaseAuth.DefaultInstance.SetCustomUserClaimsAsync(uid, customClaims);
            _logger.LogInformation("Set custom claims for user: {Uid}", uid);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Failed to set custom claims for user: {Uid}", uid);
            throw;
        }
    }

    // ==================== 사용자 목록 ====================

    public async Task<IEnumerable<User>> GetAllUsersAsync()
    {
        // PostgreSQL에서 모든 사용자 조회
        return await Task.FromResult(_userRepository.GetAllAsync().Result);
    }
}
