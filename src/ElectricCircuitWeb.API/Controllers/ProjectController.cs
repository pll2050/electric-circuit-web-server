using Microsoft.AspNetCore.Mvc;
using ElectricCircuitWeb.API.Services;
using ElectricCircuitWeb.API.Models;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 프로젝트 컨트롤러
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class ProjectController : ControllerBase
{
    private readonly IProjectService _projectService;
    private readonly ILogger<ProjectController> _logger;

    public ProjectController(IProjectService projectService, ILogger<ProjectController> logger)
    {
        _projectService = projectService;
        _logger = logger;
    }

    // ==================== 프로젝트 목록 조회 ====================

    /// <summary>
    /// 사용자의 모든 프로젝트를 조회합니다.
    /// GET /api/projects
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> GetProjects()
    {
        try
        {
            // TODO: 실제 환경에서는 인증 미들웨어에서 사용자 ID 추출
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            var projects = await _projectService.GetUserProjectsAsync(userId);

            return Ok(new
            {
                success = true,
                message = "Projects retrieved successfully",
                projects = projects.Select(p => new
                {
                    id = p.Id,
                    name = p.Name,
                    description = p.Description,
                    user_id = p.OwnerId,
                    status = "active", // 기본값
                    created_at = p.CreatedAt,
                    updated_at = p.UpdatedAt
                })
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting projects");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 프로젝트 생성 ====================

    /// <summary>
    /// 새로운 프로젝트를 생성합니다.
    /// POST /api/projects/create
    /// </summary>
    [HttpPost("create")]
    public async Task<IActionResult> CreateProject([FromBody] CreateProjectRequest request)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(request.Name))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Project name is required"
                });
            }

            var project = new Project
            {
                Name = request.Name,
                Description = request.Description ?? string.Empty,
                OwnerId = userId
            };

            var createdProject = await _projectService.CreateProjectAsync(project);

            return Ok(new
            {
                success = true,
                message = "Project created successfully",
                project = new
                {
                    id = createdProject.Id,
                    name = createdProject.Name,
                    description = createdProject.Description,
                    user_id = createdProject.OwnerId,
                    status = "active",
                    created_at = createdProject.CreatedAt,
                    updated_at = createdProject.UpdatedAt
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating project");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 프로젝트 상세 조회 ====================

    /// <summary>
    /// 특정 프로젝트의 상세 정보를 조회합니다.
    /// GET /api/projects/get?projectId={project_id}
    /// </summary>
    [HttpGet("get")]
    public async Task<IActionResult> GetProject([FromQuery] string projectId)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(projectId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Project ID is required"
                });
            }

            var project = await _projectService.GetProjectByIdAsync(projectId);

            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            // 권한 확인: 프로젝트 소유자만 조회 가능
            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            return Ok(new
            {
                success = true,
                message = "Project retrieved successfully",
                project = new
                {
                    id = project.Id,
                    name = project.Name,
                    description = project.Description,
                    user_id = project.OwnerId,
                    status = "active",
                    settings = new
                    {
                        grid_size = 10,
                        snap_to_grid = true
                    },
                    created_at = project.CreatedAt,
                    updated_at = project.UpdatedAt
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting project: {ProjectId}", projectId);
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 프로젝트 수정 ====================

    /// <summary>
    /// 프로젝트 정보를 수정합니다.
    /// PUT /api/projects/update
    /// </summary>
    [HttpPut("update")]
    public async Task<IActionResult> UpdateProject([FromBody] UpdateProjectRequest request)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(request.ProjectId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Project ID is required"
                });
            }

            var project = await _projectService.GetProjectByIdAsync(request.ProjectId);

            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            // 권한 확인
            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            // 업데이트
            if (!string.IsNullOrEmpty(request.Name))
                project.Name = request.Name;
            if (!string.IsNullOrEmpty(request.Description))
                project.Description = request.Description;

            var updatedProject = await _projectService.UpdateProjectAsync(project);

            return Ok(new
            {
                success = true,
                message = "Project updated successfully",
                project = new
                {
                    id = updatedProject.Id,
                    name = updatedProject.Name,
                    description = updatedProject.Description,
                    updated_at = updatedProject.UpdatedAt
                }
            });
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new
            {
                success = false,
                error = ex.Message
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error updating project: {ProjectId}", request.ProjectId);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 프로젝트 삭제 ====================

    /// <summary>
    /// 프로젝트를 삭제합니다.
    /// DELETE /api/projects/delete?projectId={project_id}
    /// </summary>
    [HttpDelete("delete")]
    public async Task<IActionResult> DeleteProject([FromQuery] string projectId)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(projectId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Project ID is required"
                });
            }

            var project = await _projectService.GetProjectByIdAsync(projectId);

            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            // 권한 확인
            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            var deleted = await _projectService.DeleteProjectAsync(projectId);

            if (!deleted)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            return Ok(new
            {
                success = true,
                message = "Project deleted successfully"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting project: {ProjectId}", projectId);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 프로젝트 복제 ====================

    /// <summary>
    /// 기존 프로젝트를 복제합니다.
    /// POST /api/projects/duplicate
    /// </summary>
    [HttpPost("duplicate")]
    public async Task<IActionResult> DuplicateProject([FromBody] DuplicateProjectRequest request)
    {
        try
        {
            var userId = HttpContext.Request.Headers["X-User-ID"].FirstOrDefault();

            if (string.IsNullOrEmpty(userId))
            {
                return Unauthorized(new
                {
                    success = false,
                    error = "User not authenticated"
                });
            }

            if (string.IsNullOrEmpty(request.ProjectId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Project ID is required"
                });
            }

            if (string.IsNullOrEmpty(request.Name))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "New project name is required"
                });
            }

            var originalProject = await _projectService.GetProjectByIdAsync(request.ProjectId);

            if (originalProject == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Original project not found"
                });
            }

            // 권한 확인
            if (originalProject.OwnerId != userId)
            {
                return Forbid();
            }

            var duplicatedProject = await _projectService.DuplicateProjectAsync(request.ProjectId, request.Name);

            return Ok(new
            {
                success = true,
                message = "Project duplicated successfully",
                project = new
                {
                    id = duplicatedProject.Id,
                    name = duplicatedProject.Name,
                    description = duplicatedProject.Description,
                    user_id = duplicatedProject.OwnerId,
                    status = "active",
                    created_at = duplicatedProject.CreatedAt
                }
            });
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new
            {
                success = false,
                error = ex.Message
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error duplicating project: {ProjectId}", request.ProjectId);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }
}

// ==================== Request Models ====================

public record CreateProjectRequest(string Name, string? Description = null);
public record UpdateProjectRequest(string ProjectId, string? Name = null, string? Description = null);
public record DuplicateProjectRequest(string ProjectId, string Name);
