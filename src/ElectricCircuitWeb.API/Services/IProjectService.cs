using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 프로젝트 서비스 인터페이스
/// </summary>
public interface IProjectService
{
    // 프로젝트 조회
    Task<IEnumerable<Project>> GetUserProjectsAsync(string userId);
    Task<Project?> GetProjectByIdAsync(string projectId);

    // 프로젝트 생성
    Task<Project> CreateProjectAsync(Project project);

    // 프로젝트 수정
    Task<Project> UpdateProjectAsync(Project project);

    // 프로젝트 삭제
    Task<bool> DeleteProjectAsync(string projectId);

    // 프로젝트 복제
    Task<Project> DuplicateProjectAsync(string projectId, string newName);
}
