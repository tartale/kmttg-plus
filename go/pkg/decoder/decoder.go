package decoder

import (
	"context"
	"io"
	"os/exec"

	"github.com/tartale/kmttg-plus/go/pkg/config"
	"github.com/tartale/kmttg-plus/go/pkg/logz"
)

func Decode(ctx context.Context, in io.Reader, out io.Writer) error {
	decoderCommand := config.TivoDecoderCmd
	decoder := exec.CommandContext(ctx, decoderCommand[0], decoderCommand[1:]...)
	decoder.Stdin = in
	decoder.Stdout = out
	logz.LoggerX.Debugf("Start decoding")
	if err := decoder.Run(); err != nil {
		logz.LoggerX.Errorf("%w: error running decoder", err)
		return err
	}
	logz.LoggerX.Debugf("Finished decoding")
	return nil
}
