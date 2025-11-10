using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Repositories;

/// <summary>
/// 프로젝트 리포지토리 인터페이스
/// </summary>
public interface IProjectRepository
{
    Task<IEnumerable<Project>> GetByUserIdAsync(string userId);
    Task<Project?> GetByIdAsync(string projectId);
    Task<Project> CreateAsync(Project project);
    Task<Project> UpdateAsync(Project project);
    Task<bool> DeleteAsync(string projectId);
}
