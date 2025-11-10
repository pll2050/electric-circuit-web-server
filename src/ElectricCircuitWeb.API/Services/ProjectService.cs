using ElectricCircuitWeb.API.Models;
using ElectricCircuitWeb.API.Repositories;

namespace ElectricCircuitWeb.API.Services;

/// <summary>
/// 프로젝트 서비스 구현
/// </summary>
public class ProjectService : IProjectService
{
    private readonly IProjectRepository _projectRepository;
    private readonly ILogger<ProjectService> _logger;

    public ProjectService(IProjectRepository projectRepository, ILogger<ProjectService> logger)
    {
        _projectRepository = projectRepository;
        _logger = logger;
    }

    public async Task<IEnumerable<Project>> GetUserProjectsAsync(string userId)
    {
        return await _projectRepository.GetByUserIdAsync(userId);
    }

    public async Task<Project?> GetProjectByIdAsync(string projectId)
    {
        return await _projectRepository.GetByIdAsync(projectId);
    }

    public async Task<Project> CreateProjectAsync(Project project)
    {
        // ID가 없으면 생성
        if (string.IsNullOrEmpty(project.Id))
        {
            project.Id = Guid.NewGuid().ToString();
        }

        project.CreatedAt = DateTime.UtcNow;
        project.UpdatedAt = null;

        return await _projectRepository.CreateAsync(project);
    }

    public async Task<Project> UpdateProjectAsync(Project project)
    {
        var existingProject = await _projectRepository.GetByIdAsync(project.Id);
        if (existingProject == null)
        {
            throw new KeyNotFoundException($"Project with ID {project.Id} not found");
        }

        return await _projectRepository.UpdateAsync(project);
    }

    public async Task<bool> DeleteProjectAsync(string projectId)
    {
        return await _projectRepository.DeleteAsync(projectId);
    }

    public async Task<Project> DuplicateProjectAsync(string projectId, string newName)
    {
        var originalProject = await _projectRepository.GetByIdAsync(projectId);
        if (originalProject == null)
        {
            throw new KeyNotFoundException($"Project with ID {projectId} not found");
        }

        var duplicatedProject = new Project
        {
            Id = Guid.NewGuid().ToString(),
            Name = newName,
            Description = originalProject.Description,
            OwnerId = originalProject.OwnerId,
            CreatedAt = DateTime.UtcNow,
            UpdatedAt = null
        };

        return await _projectRepository.CreateAsync(duplicatedProject);
    }
}
