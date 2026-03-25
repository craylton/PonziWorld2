import { useEffect, useState } from 'react';
import PonziFactorSlider from './PonziValue';
import { makeAuthenticatedRequest } from '../../auth';

interface TomorrowPanelProps {
  bankId: string;
}

export default function TomorrowPanel({ bankId }: TomorrowPanelProps) {
  const [ponziFactor, setPonziFactor] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [hasLoadError, setHasLoadError] = useState(false);

  useEffect(() => {
    const abortController = new AbortController();

    const load = async () => {
      setIsLoading(true);
      setHasLoadError(false);

      try {
        const response = await makeAuthenticatedRequest(
          `/api/bank/ponziFactor?bankId=${encodeURIComponent(bankId)}`,
          { signal: abortController.signal },
        );

        if (!response.ok) {
          throw new Error('Load failed');
        }

        const data = (await response.json()) as { ponziFactor: number };
        setPonziFactor(data.ponziFactor);
      } catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError') {
          return;
        }

        setHasLoadError(true);
      } finally {
        if (!abortController.signal.aborted) {
          setIsLoading(false);
        }
      }
    };

    void load();

    return () => {
      abortController.abort();
    };
  }, [bankId]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (hasLoadError || ponziFactor === null) {
    return <div>Error loading ponzi factor</div>;
  }

  return (
    <PonziFactorSlider
      bankId={bankId}
      ponziFactor={ponziFactor}
      onPonziFactorSaved={setPonziFactor}
    />
  );
}
