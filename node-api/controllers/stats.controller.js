// Calcula estadísticas sobre las matrices Q y R
export function calculateStats(req, res) {
    const { Q, R } = req.body;
    if (!Array.isArray(Q) || !Array.isArray(R)) {
        return res.status(400).json({ error: 'Q y R deben ser matrices' });
    }
  
    const allValues = [...Q.flat(), ...R.flat()];
    const max = Math.max(...allValues);
    const min = Math.min(...allValues);
    const sum = allValues.reduce((acc, val) => acc + val, 0);
    const average = sum / allValues.length;
  
    const isDiagonal = (matrix) =>
      matrix.every((row, i) =>
        row.every((val, j) => (i !== j && val !== 0) ? false : true)
      );
  
    const anyDiagonal = isDiagonal(Q) || isDiagonal(R);
  
    res.status(200).json({
      max,
      min,
      sum,
      average,
      isDiagonal: anyDiagonal
    });
  }