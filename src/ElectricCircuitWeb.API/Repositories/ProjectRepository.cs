using ElectricCircuitWeb.API.Data;
using ElectricCircuitWeb.API.Models;
using Microsoft.EntityFrameworkCore;

namespace ElectricCircuitWeb.API.Repositories;

/// <summary>
/// 프로젝트 리포지토리 구현
/// </summary>
public class ProjectRepository : IProjectRepository
{
    private readonly ApplicationDbContext _context;

    public ProjectRepository(ApplicationDbContext context)
    {
        _context = context;
    }

    public async Task<IEnumerable<Project>> GetByUserIdAsync(string userId)
    {
        return await _context.Projects
            .Where(p => p.OwnerId == userId)
            .OrderByDescending(p => p.CreatedAt)
            .ToListAsync();
    }

    public async Task<Project?> GetByIdAsync(string projectId)
    {
        return await _context.Projects
            .FirstOrDefaultAsync(p => p.Id == projectId);
    }

    public async Task<Project> CreateAsync(Project project)
    {
        _context.Projects.Add(project);
        await _context.SaveChangesAsync();
        return project;
    }

    public async Task<Project> UpdateAsync(Project project)
    {
        project.UpdatedAt = DateTime.UtcNow;
        _context.Projects.Update(project);
        await _context.SaveChangesAsync();
        return project;
    }

    public async Task<bool> DeleteAsync(string projectId)
    {
        var project = await _context.Projects.FindAsync(projectId);
        if (project == null)
            return false;

        _context.Projects.Remove(project);
        await _context.SaveChangesAsync();
        return true;
    }
}
