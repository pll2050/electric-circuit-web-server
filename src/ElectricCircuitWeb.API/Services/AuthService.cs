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

    public async Task<User?> GetUserByFirebaseUidAsync(string firebaseUid)
    {
        return await _userRepository.GetByFirebaseUidAsync(firebaseUid);
    }

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
}
