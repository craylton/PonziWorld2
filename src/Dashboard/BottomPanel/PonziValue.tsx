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
  const sanitizedInitialPonziFactor = sanitizePonziFactor(ponziFactor);

  const [ponziFactorInput, setPonziFactorInput] = useState<number>(() => sanitizedInitialPonziFactor);
  const [saveStatus, setSaveStatus] = useState<SaveStatus>('idle');

  const isMountedRef = useRef(false);
  const ponziFactorInputRef = useRef<number>(sanitizedInitialPonziFactor);
  const lastSavedPonziFactorRef = useRef<number>(sanitizedInitialPonziFactor);

  const debounceTimeoutIdRef = useRef<number | null>(null);
  const clearSavedMessageTimeoutIdRef = useRef<number | null>(null);
  const saveAbortControllerRef = useRef<AbortController | null>(null);

  const clearDebounceTimeout = useCallback(() => {
    if (debounceTimeoutIdRef.current !== null) {
      window.clearTimeout(debounceTimeoutIdRef.current);
      debounceTimeoutIdRef.current = null;
    }
  }, []);

  const clearSavedMessageTimeout = useCallback(() => {
    if (clearSavedMessageTimeoutIdRef.current !== null) {
      window.clearTimeout(clearSavedMessageTimeoutIdRef.current);
      clearSavedMessageTimeoutIdRef.current = null;
    }
  }, []);

  const markSaving = useCallback(() => {
    clearSavedMessageTimeout();
    if (isMountedRef.current) {
      setSaveStatus('saving');
    }
  }, [clearSavedMessageTimeout]);

  const scheduleClearSavedMessage = useCallback(() => {
    clearSavedMessageTimeout();
    clearSavedMessageTimeoutIdRef.current = window.setTimeout(() => {
      if (!isMountedRef.current) {
        return;
      }

      setSaveStatus('idle');
    }, savedMessageMilliseconds);
  }, [clearSavedMessageTimeout]);

  const savePonziFactor = useCallback(async (valueToSave: number) => {
    const sanitizedValueToSave = sanitizePonziFactor(valueToSave);

    if (sanitizedValueToSave === lastSavedPonziFactorRef.current) {
      return;
    }

    if (saveAbortControllerRef.current) {
      saveAbortControllerRef.current.abort();
    }

    const abortController = new AbortController();
    saveAbortControllerRef.current = abortController;

    markSaving();

    try {
      const response = await makeAuthenticatedRequest('/api/bank/ponziFactor', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        signal: abortController.signal,
        body: JSON.stringify({ bankId, ponziFactor: sanitizedValueToSave }),
      });

      if (saveAbortControllerRef.current !== abortController) {
        return;
      }

      if (!response.ok) {
        throw new Error('Save failed');
      }

      lastSavedPonziFactorRef.current = sanitizedValueToSave;

      onPonziFactorSaved(sanitizedValueToSave);

      if (!isMountedRef.current) {
        return;
      }

      setSaveStatus('saved');

      scheduleClearSavedMessage();
    } catch (error) {
      if (saveAbortControllerRef.current !== abortController) {
        return;
      }

      if (error instanceof DOMException && error.name === 'AbortError') {
        return;
      }

      if (!isMountedRef.current) {
        return;
      }

      clearSavedMessageTimeout();
      setSaveStatus('error');
    } finally {
      if (saveAbortControllerRef.current === abortController) {
        saveAbortControllerRef.current = null;
      }
    }
  }, [bankId, clearSavedMessageTimeout, markSaving, onPonziFactorSaved, scheduleClearSavedMessage]);

  const scheduleDebouncedSave = useCallback((valueToSave: number) => {
    clearDebounceTimeout();
    debounceTimeoutIdRef.current = window.setTimeout(() => {
      debounceTimeoutIdRef.current = null;
      void savePonziFactor(valueToSave);
    }, debounceMilliseconds);
  }, [clearDebounceTimeout, savePonziFactor]);

  useEffect(() => {
    isMountedRef.current = true;
    return () => {
      isMountedRef.current = false;

      const shouldFlushDebouncedSave = debounceTimeoutIdRef.current !== null;
      clearDebounceTimeout();
      clearSavedMessageTimeout();

      if (!shouldFlushDebouncedSave) {
        return;
      }

      const latestValue = ponziFactorInputRef.current;
      if (latestValue === lastSavedPonziFactorRef.current) {
        return;
      }

      void savePonziFactor(latestValue);
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
    ponziFactorInputRef.current = sanitizedNextValue;
    setPonziFactorInput(sanitizedNextValue);

    if (sanitizedNextValue === lastSavedPonziFactorRef.current) {
      clearDebounceTimeout();
      clearSavedMessageTimeout();
      setSaveStatus('idle');
      return;
    }

    markSaving();
    scheduleDebouncedSave(sanitizedNextValue);
  };

  const statusText = saveStatus === 'saving' ? 'Saving...' : saveStatus === 'saved' ? 'Saved' : saveStatus === 'error' ? 'Error' : '';
  const statusClassName = saveStatus === 'saving'
    ? 'ponzi-value__status ponzi-value__status--saving'
    : saveStatus === 'saved'
      ? 'ponzi-value__status ponzi-value__status--saved'
      : saveStatus === 'error'
        ? 'ponzi-value__status ponzi-value__status--error'
        : 'ponzi-value__status';

  return (
    <div className="ponzi-value">
      <label className="ponzi-value__label">
        Ponzi value:{' '}
        <span
          className="ponzi-value__percent"
          style={{
            display: 'inline-block',
            minWidth: '6ch',
            textAlign: 'right',
            fontVariantNumeric: 'tabular-nums',
          }}
        >
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
        <span className={statusClassName} aria-live="polite">
          {statusText}
        </span>
      </div>
    </div>
  );
}
