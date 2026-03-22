import { useCallback, useEffect, useRef, useState } from 'react';
import { makeAuthenticatedRequest } from '../../auth';

interface PonziValueProps {
  bankId: string;
  ponziFactor: number;
  onPonziFactorSaved: (ponziFactor: number) => void;
}

const debounceMilliseconds = 500;
const savedMessageMilliseconds = 2000;

const minPonziFactor = -0.05;
const maxPonziFactor = 0.15;

type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';

const statusText: Record<SaveStatus, string> = {
  idle: '',
  saving: 'Saving...',
  saved: 'Saved',
  error: 'Error',
};

const sanitizePonziFactor = (value: number) => {
  if (!Number.isFinite(value)) {
    return 0;
  }

  return Math.max(minPonziFactor, Math.min(maxPonziFactor, value));
};

const formatPercent = (value: number) => {
  const percent = value * 100;
  const rounded = Math.round(percent * 10) / 10;

  if (Number.isInteger(rounded)) {
    return `${rounded.toFixed(0)}%`;
  }

  return `${rounded.toFixed(1)}%`;
};

export default function PonziValue({ bankId, ponziFactor, onPonziFactorSaved }: PonziValueProps) {
  const sanitizedPonziFactor = sanitizePonziFactor(ponziFactor);

  const [ponziFactorInput, setPonziFactorInput] = useState(sanitizedPonziFactor);
  const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle');

  const ponziFactorInputRef = useRef<number>(sanitizedPonziFactor);
  const lastSavedPonziFactorRef = useRef<number>(sanitizedPonziFactor);

  const debounceTimeoutIdRef = useRef<number | null>(null);
  const savedMessageTimeoutIdRef = useRef<number | null>(null);
  const saveAbortControllerRef = useRef<AbortController | null>(null);

  const clearDebounceTimeout = useCallback(() => {
    if (debounceTimeoutIdRef.current !== null) {
      window.clearTimeout(debounceTimeoutIdRef.current);
      debounceTimeoutIdRef.current = null;
    }
  }, []);

  const clearSavedMessageTimeout = useCallback(() => {
    if (savedMessageTimeoutIdRef.current !== null) {
      window.clearTimeout(savedMessageTimeoutIdRef.current);
      savedMessageTimeoutIdRef.current = null;
    }
  }, []);

  const updateSaveStatus = useCallback((nextStatus: SaveStatus) => {
    clearSavedMessageTimeout();
    setSaveStatus(nextStatus);
  }, [clearSavedMessageTimeout]);

  const showSavedStatus = useCallback(() => {
    updateSaveStatus('saved');
    savedMessageTimeoutIdRef.current = window.setTimeout(() => {
      savedMessageTimeoutIdRef.current = null;
      setSaveStatus('idle');
    }, savedMessageMilliseconds);
  }, [updateSaveStatus]);

  const savePonziFactor = useCallback(async (sanitizedValue: number) => {
    if (sanitizedValue === lastSavedPonziFactorRef.current) {
      return;
    }

    saveAbortControllerRef.current?.abort();

    const abortController = new AbortController();
    saveAbortControllerRef.current = abortController;

    updateSaveStatus('saving');

    try {
      const response = await makeAuthenticatedRequest('/api/bank/ponziFactor', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        signal: abortController.signal,
        body: JSON.stringify({ bankId, ponziFactor: sanitizedValue }),
      });

      if (saveAbortControllerRef.current !== abortController) {
        return;
      }

      if (!response.ok) {
        throw new Error('Save failed');
      }

      lastSavedPonziFactorRef.current = sanitizedValue;

      onPonziFactorSaved(sanitizedValue);

      showSavedStatus();
    } catch (error) {
      if (saveAbortControllerRef.current !== abortController) {
        return;
      }

      if (error instanceof DOMException && error.name === 'AbortError') {
        return;
      }

      updateSaveStatus('error');
    } finally {
      if (saveAbortControllerRef.current === abortController) {
        saveAbortControllerRef.current = null;
      }
    }
  }, [bankId, onPonziFactorSaved, showSavedStatus, updateSaveStatus]);

  const scheduleDebouncedSave = useCallback((valueToSave: number) => {
    clearDebounceTimeout();
    debounceTimeoutIdRef.current = window.setTimeout(() => {
      debounceTimeoutIdRef.current = null;
      void savePonziFactor(valueToSave);
    }, debounceMilliseconds);
  }, [clearDebounceTimeout, savePonziFactor]);

  useEffect(() => {
    return () => {
      const shouldFlushDebouncedSave = debounceTimeoutIdRef.current !== null;
      clearDebounceTimeout();
      clearSavedMessageTimeout();

      if (!shouldFlushDebouncedSave) {
        return;
      }

      void savePonziFactor(ponziFactorInputRef.current);
    };
  }, [clearDebounceTimeout, clearSavedMessageTimeout, savePonziFactor]);

  useEffect(() => {
    const sanitizedPropValue = sanitizePonziFactor(ponziFactor);

    const hasUnsavedLocalChanges = ponziFactorInputRef.current !== lastSavedPonziFactorRef.current;
    if (hasUnsavedLocalChanges) {
      return;
    }

    ponziFactorInputRef.current = sanitizedPropValue;
    lastSavedPonziFactorRef.current = sanitizedPropValue;
    setPonziFactorInput(sanitizedPropValue);
  }, [ponziFactor]);

  const handleChange = (nextValue: number) => {
    const sanitizedNextValue = sanitizePonziFactor(nextValue);

    saveAbortControllerRef.current?.abort();

    ponziFactorInputRef.current = sanitizedNextValue;
    setPonziFactorInput(sanitizedNextValue);

    if (sanitizedNextValue === lastSavedPonziFactorRef.current) {
      clearDebounceTimeout();
      updateSaveStatus('idle');
      return;
    }

    updateSaveStatus('saving');
    scheduleDebouncedSave(sanitizedNextValue);
  };

  const statusModifier = saveStatus !== 'idle' ? ` ponzi-value__status--${saveStatus}` : '';

  return (
    <div className="ponzi-value">
      <label className="ponzi-value__label">
        Ponzi value:{' '}
        <span className="ponzi-value__percent">
          {formatPercent(ponziFactorInput)}
        </span>
      </label>

      <div className="ponzi-value__controls">
        <input
          type="range"
          min={minPonziFactor}
          max={maxPonziFactor}
          step={0.001}
          value={ponziFactorInput}
          onChange={(event) => handleChange(Number(event.target.value))}
        />
        <span className={`ponzi-value__status${statusModifier}`} aria-live="polite">
          {statusText[saveStatus]}
        </span>
      </div>
    </div>
  );
}
