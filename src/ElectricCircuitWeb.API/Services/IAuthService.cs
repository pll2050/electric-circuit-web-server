using ElectricCircuitWeb.API.Models;
using FirebaseAdmin.Auth;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 인증 서비스 인터페이스
/// </summary>
public interface IAuthService
{
    // 토큰 검증
    Task<string> VerifyIdTokenAsync(string idToken);

    // 사용자 조회
    Task<User?> GetUserByFirebaseUidAsync(string firebaseUid);
    Task<UserRecord?> GetFirebaseUserAsync(string uid);

    // 사용자 생성
    Task<User> CreateUserAsync(string firebaseUid, string email, string displayName);
    Task<UserRecord> CreateFirebaseUserAsync(string email, string password, string? displayName = null, string? photoUrl = null);

    // 사용자 수정
    Task<User> UpdateUserAsync(User user);
    Task<UserRecord> UpdateFirebaseUserAsync(string uid, string? email = null, string? displayName = null, string? photoUrl = null, string? password = null);

    // 사용자 삭제
    Task<bool> DeleteUserAsync(string firebaseUid);
    Task DeleteFirebaseUserAsync(string uid);

    // 커스텀 클레임
    Task SetCustomClaimsAsync(string uid, Dictionary<string, object> customClaims);

    // 사용자 목록
    Task<IEnumerable<User>> GetAllUsersAsync();
}
