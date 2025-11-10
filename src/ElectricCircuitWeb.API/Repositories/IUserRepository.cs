using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Repositories;

/// <summary>
/// 사용자 리포지토리 인터페이스
/// </summary>
public interface IUserRepository
{
    Task<User?> GetByIdAsync(int id);
    Task<User?> GetByFirebaseUidAsync(string firebaseUid);
    Task<IEnumerable<User>> GetAllAsync();
    Task<User> CreateAsync(User user);
    Task<User> UpdateAsync(User user);
    Task<bool> DeleteAsync(int id);
}
