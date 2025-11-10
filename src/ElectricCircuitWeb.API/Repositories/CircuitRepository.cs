using ElectricCircuitWeb.API.Data;
using ElectricCircuitWeb.API.Models;
using Microsoft.EntityFrameworkCore;

namespace ElectricCircuitWeb.API.Repositories;

/// <summary>
/// 회로 리포지토리 구현
/// </summary>
public class CircuitRepository : ICircuitRepository
{
    private readonly ApplicationDbContext _context;

    public CircuitRepository(ApplicationDbContext context)
    {
        _context = context;
    }

    public async Task<IEnumerable<Circuit>> GetByProjectIdAsync(string projectId)
    {
        return await _context.Circuits
            .Where(c => c.ProjectId == projectId)
            .OrderByDescending(c => c.CreatedAt)
            .ToListAsync();
    }

    public async Task<Circuit?> GetByIdAsync(string circuitId)
    {
        return await _context.Circuits
            .FirstOrDefaultAsync(c => c.Id == circuitId);
    }

    public async Task<Circuit> CreateAsync(Circuit circuit)
    {
        _context.Circuits.Add(circuit);
        await _context.SaveChangesAsync();
        return circuit;
    }

    public async Task<Circuit> UpdateAsync(Circuit circuit)
    {
        circuit.UpdatedAt = DateTime.UtcNow;
        _context.Circuits.Update(circuit);
        await _context.SaveChangesAsync();
        return circuit;
    }

    public async Task<bool> DeleteAsync(string circuitId)
    {
        var circuit = await _context.Circuits
            .FirstOrDefaultAsync(c => c.Id == circuitId);

        if (circuit == null)
            return false;

        _context.Circuits.Remove(circuit);
        await _context.SaveChangesAsync();
        return true;
    }
}
