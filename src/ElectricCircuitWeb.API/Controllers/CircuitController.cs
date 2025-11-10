using Microsoft.AspNetCore.Mvc;
using ElectricCircuitWeb.API.Services;
using ElectricCircuitWeb.API.Models;
using System.Text.Json;

namespace ElectricCircuitWeb.API.Controllers;

/// <summary>
/// 회로 컨트롤러
/// </summary>
[ApiController]
[Route("api/[controller]")]
public class CircuitController : ControllerBase
{
    private readonly ICircuitService _circuitService;
    private readonly IProjectService _projectService;
    private readonly ILogger<CircuitController> _logger;

    public CircuitController(
        ICircuitService circuitService,
        IProjectService projectService,
        ILogger<CircuitController> logger)
    {
        _circuitService = circuitService;
        _projectService = projectService;
        _logger = logger;
    }

    // ==================== 회로 목록 조회 ====================

    /// <summary>
    /// 특정 프로젝트의 모든 회로를 조회합니다.
    /// GET /api/circuits?projectId={project_id}
    /// </summary>
    [HttpGet]
    public async Task<IActionResult> GetCircuits([FromQuery] string projectId)
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

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(projectId);
            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            var circuits = await _circuitService.GetProjectCircuitsAsync(projectId);

            return Ok(new
            {
                success = true,
                message = "Circuits retrieved successfully",
                circuits = circuits.Select(c => new
                {
                    id = c.Id,
                    name = c.Name,
                    description = "",
                    project_id = c.ProjectId,
                    user_id = userId,
                    version = 1,
                    is_template = false,
                    tags = new string[] { },
                    created_at = c.CreatedAt,
                    updated_at = c.UpdatedAt
                })
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting circuits for project: {ProjectId}", projectId);
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 회로 생성 ====================

    /// <summary>
    /// 새로운 회로를 생성합니다.
    /// POST /api/circuits/create
    /// </summary>
    [HttpPost("create")]
    public async Task<IActionResult> CreateCircuit([FromBody] CreateCircuitRequest request)
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
                    error = "Circuit name is required"
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

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(request.ProjectId);
            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            var circuit = new Circuit
            {
                Name = request.Name,
                ProjectId = request.ProjectId,
                Data = request.Data != null ? JsonSerializer.Serialize(request.Data) : "{}"
            };

            var createdCircuit = await _circuitService.CreateCircuitAsync(circuit);

            return Ok(new
            {
                success = true,
                message = "Circuit created successfully",
                circuit = new
                {
                    id = createdCircuit.Id,
                    name = createdCircuit.Name,
                    description = "",
                    project_id = createdCircuit.ProjectId,
                    user_id = userId,
                    data = JsonSerializer.Deserialize<object>(createdCircuit.Data),
                    version = 1,
                    created_at = createdCircuit.CreatedAt
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating circuit");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 회로 상세 조회 ====================

    /// <summary>
    /// 특정 회로의 상세 정보를 조회합니다.
    /// GET /api/circuits/get?circuitId={circuit_id}
    /// </summary>
    [HttpGet("get")]
    public async Task<IActionResult> GetCircuit([FromQuery] string circuitId)
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

            if (string.IsNullOrEmpty(circuitId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Circuit ID is required"
                });
            }

            var circuit = await _circuitService.GetCircuitByIdAsync(circuitId);

            if (circuit == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Circuit not found"
                });
            }

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(circuit.ProjectId);
            if (project == null || project.OwnerId != userId)
            {
                return Forbid();
            }

            return Ok(new
            {
                success = true,
                message = "Circuit retrieved successfully",
                circuit = new
                {
                    id = circuit.Id,
                    name = circuit.Name,
                    description = "",
                    project_id = circuit.ProjectId,
                    user_id = userId,
                    data = JsonSerializer.Deserialize<object>(circuit.Data),
                    version = 1,
                    tags = new string[] { },
                    created_at = circuit.CreatedAt,
                    updated_at = circuit.UpdatedAt
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting circuit: {CircuitId}", circuitId);
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    // ==================== 회로 수정 ====================

    /// <summary>
    /// 회로 정보 및 설계 데이터를 수정합니다.
    /// PUT /api/circuits/update
    /// </summary>
    [HttpPut("update")]
    public async Task<IActionResult> UpdateCircuit([FromBody] UpdateCircuitRequest request)
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

            if (string.IsNullOrEmpty(request.CircuitId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Circuit ID is required"
                });
            }

            var circuit = await _circuitService.GetCircuitByIdAsync(request.CircuitId);

            if (circuit == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Circuit not found"
                });
            }

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(circuit.ProjectId);
            if (project == null || project.OwnerId != userId)
            {
                return Forbid();
            }

            // 업데이트
            if (!string.IsNullOrEmpty(request.Name))
                circuit.Name = request.Name;
            if (request.Data != null)
                circuit.Data = JsonSerializer.Serialize(request.Data);

            var updatedCircuit = await _circuitService.UpdateCircuitAsync(circuit);

            return Ok(new
            {
                success = true,
                message = "Circuit updated successfully",
                circuit = new
                {
                    id = updatedCircuit.Id,
                    name = updatedCircuit.Name,
                    description = "",
                    version = 2,
                    updated_at = updatedCircuit.UpdatedAt
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
            _logger.LogError(ex, "Error updating circuit: {CircuitId}", request.CircuitId);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 회로 삭제 ====================

    /// <summary>
    /// 회로를 삭제합니다.
    /// DELETE /api/circuits/delete?circuitId={circuit_id}
    /// </summary>
    [HttpDelete("delete")]
    public async Task<IActionResult> DeleteCircuit([FromQuery] string circuitId)
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

            if (string.IsNullOrEmpty(circuitId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Circuit ID is required"
                });
            }

            var circuit = await _circuitService.GetCircuitByIdAsync(circuitId);

            if (circuit == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Circuit not found"
                });
            }

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(circuit.ProjectId);
            if (project == null || project.OwnerId != userId)
            {
                return Forbid();
            }

            var deleted = await _circuitService.DeleteCircuitAsync(circuitId);

            if (!deleted)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Circuit not found"
                });
            }

            return Ok(new
            {
                success = true,
                message = "Circuit deleted successfully"
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error deleting circuit: {CircuitId}", circuitId);
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }

    // ==================== 회로 템플릿 ====================

    /// <summary>
    /// 사용 가능한 회로 템플릿을 조회합니다.
    /// GET /api/circuits/templates
    /// </summary>
    [HttpGet("templates")]
    public async Task<IActionResult> GetTemplates()
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

            var templates = await _circuitService.GetTemplatesAsync();

            return Ok(new
            {
                success = true,
                message = "Templates feature not yet implemented",
                templates
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error getting templates");
            return StatusCode(500, new
            {
                success = false,
                error = "Internal server error"
            });
        }
    }

    /// <summary>
    /// 기존 템플릿을 사용하여 새 회로를 생성합니다.
    /// POST /api/circuits/create-from-template
    /// </summary>
    [HttpPost("create-from-template")]
    public async Task<IActionResult> CreateFromTemplate([FromBody] CreateFromTemplateRequest request)
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

            if (string.IsNullOrEmpty(request.TemplateId))
            {
                return BadRequest(new
                {
                    success = false,
                    error = "Template ID is required"
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
                    error = "Circuit name is required"
                });
            }

            // 프로젝트 권한 확인
            var project = await _projectService.GetProjectByIdAsync(request.ProjectId);
            if (project == null)
            {
                return NotFound(new
                {
                    success = false,
                    error = "Project not found"
                });
            }

            if (project.OwnerId != userId)
            {
                return Forbid();
            }

            var circuit = await _circuitService.CreateFromTemplateAsync(
                request.TemplateId,
                request.ProjectId,
                request.Name
            );

            return Ok(new
            {
                success = true,
                message = "Circuit created from template successfully",
                circuit = new
                {
                    id = circuit.Id,
                    name = circuit.Name,
                    project_id = circuit.ProjectId,
                    created_at = circuit.CreatedAt
                }
            });
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Error creating circuit from template");
            return StatusCode(500, new
            {
                success = false,
                error = ex.Message
            });
        }
    }
}

// ==================== Request Models ====================

public record CreateCircuitRequest(string ProjectId, string Name, object? Data = null);
public record UpdateCircuitRequest(string CircuitId, string? Name = null, object? Data = null);
public record CreateFromTemplateRequest(string ProjectId, string TemplateId, string Name);
