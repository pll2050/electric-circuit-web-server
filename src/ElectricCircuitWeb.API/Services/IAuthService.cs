using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 인증 서비스 인터페이스
/// </summary>
public interface IAuthService
{
    Task<User?> GetUserByFirebaseUidAsync(string firebaseUid);
    Task<User> CreateUserAsync(string firebaseUid, string email, string displayName);
    Task<string> VerifyIdTokenAsync(string idToken);
}
