// Marca do Janus (original): dois portais entrelaçados, o dourado do Fabric e o ciano do Besu.
// A interseção é a passagem guardada por Jano, deus romano das portas, com uma face voltada para
// cada mundo.
export function JanusMark({ size = 44, title = "Janus" }: { size?: number; title?: string }) {
  return (
    <svg width={size} height={size} viewBox="0 0 64 64" role="img" aria-label={title} className="janus-mark">
      <defs>
        <linearGradient id="jm-gold" x1="0" y1="0" x2="1" y2="1">
          <stop offset="0" stopColor="var(--gold-hi)" />
          <stop offset="1" stopColor="var(--gold)" />
        </linearGradient>
        <linearGradient id="jm-cyan" x1="1" y1="0" x2="0" y2="1">
          <stop offset="0" stopColor="var(--cyan-hi)" />
          <stop offset="1" stopColor="var(--cyan)" />
        </linearGradient>
        <radialGradient id="jm-core" cx="0.5" cy="0.5" r="0.5">
          <stop offset="0" stopColor="var(--violet-hi)" />
          <stop offset="1" stopColor="var(--violet)" />
        </radialGradient>
      </defs>
      {/* lente da passagem (interseção dos portais) */}
      <path d="M32 14 A20 20 0 0 1 32 50 A20 20 0 0 1 32 14 Z" fill="url(#jm-core)" opacity="0.9" />
      <circle cx="22" cy="32" r="18" fill="none" stroke="url(#jm-gold)" strokeWidth="3.2" />
      <circle cx="42" cy="32" r="18" fill="none" stroke="url(#jm-cyan)" strokeWidth="3.2" />
      {/* porta / fechadura no centro */}
      <circle cx="32" cy="28" r="3.2" fill="var(--void)" />
      <path d="M30.3 30 L33.7 30 L34.6 38 L29.4 38 Z" fill="var(--void)" />
    </svg>
  );
}
