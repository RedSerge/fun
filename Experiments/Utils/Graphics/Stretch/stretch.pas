(*	This program stretches a picture to the given size and writes
 *	the result to the output file.
 *)

program Stretch;

uses SysUtils, Types, Graphics, Interfaces;

var source, thumb : TPicture;
	width, height : Integer;

begin	
	//check if the program is called properly
	if ParamCount < 4 then begin
		writeln('Usage: stretch <original_file> <output_file> <width> <height>');
	end else begin
	//acquire parameters from the command line
		width:=StrToInt(ParamStr(3));
		height:=StrToInt(ParamStr(4));
		source:=TPicture.Create;
		source.LoadFromFile(ParamStr(1));
	//create a (supposed) thumbnail
		thumb:=TPicture.Create;
		thumb.Bitmap.SetSize(width, height);
		thumb.Bitmap.Canvas.StretchDraw(Rect(0, 0, width, height), source.Bitmap);
		source.Free;
	//export the result
		thumb.SaveToFile(ParamStr(2));
		thumb.Free;
	end;
end.