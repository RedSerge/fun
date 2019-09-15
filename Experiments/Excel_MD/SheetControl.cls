'MultiDimensional Excel - a concept of a tensor-based table, where cells may contain multiple values in various dimensions.
'Written in VBA, compiled as an Add-In for Microsoft Excel 2019 (16.x).

'This file contains the class that handles the configuration and dimension sorting process in MultiDimensional Workbook. 

'This part is provided by MS Excel VBA: class description and attributes
VERSION 1.0 CLASS
BEGIN
  MultiUse = -1  'True
END
Attribute VB_Name = "SheetControl"
Attribute VB_GlobalNameSpace = False
Attribute VB_Creatable = False
Attribute VB_PredeclaredId = False
Attribute VB_Exposed = False
'^ This part was provided by MS Excel VBA ^

'Worksheet Handler
Public WithEvents sheetevent As Worksheet

'This part is provided by MS Excel VBA: declares the module's help context displayed by the object browser as non-existent for our custom Worksheet handler.
Attribute sheetevent.VB_VarHelpID = -1
'^ This part was provided by MS Excel VBA ^

'These two procedures below rename the pages during the sorting/renaming process;
'if it's neccessary to rename current dimension #2 as #3, it's named '#3_' first;
'any duplicate found is deleted (that is, if we rename both #2 and #3 to #4 consequently,
'the last one (#3) stays, while #2 is lost.
Private Sub PreSort(ByVal index_old As Integer, ByVal index_new As Integer)
Unsorted = MDName(index_old)
Sorted = MDName(index_new) + "_"
If IsSheet(Unsorted) Then
    If IsSheet(Sorted) Then
        ActiveWorkbook.Sheets(Sorted).Visible = True
        ActiveWorkbook.Sheets(Sorted).Delete
    End If
    ActiveWorkbook.Sheets(Unsorted).Visible = True
    ActiveWorkbook.Sheets(Unsorted).Name = Sorted
End If
End Sub

'Next, we rename '#3_' from our example above to the '#3' according to the same rules.
'They, in theory, can be different - we could swap or expand dimensions, warn user, etc.
'Here, though, we simple 'overwrite' any possible duplicates.
Private Sub FinSort(ByVal index_new As Integer)
Sorted = MDName(index_new) + "_"
Restored = MDName(index_new)
If IsSheet(Sorted) Then
    If IsSheet(Restored) Then
        ActiveWorkbook.Sheets(Restored).Visible = True
        ActiveWorkbook.Sheets(Restored).Delete
    End If
    ActiveWorkbook.Sheets(Sorted).Visible = True
    ActiveWorkbook.Sheets(Sorted).Name = Restored
End If
End Sub

'The primary sorting procedure.
Private Sub SortOut(ByVal Count As Integer)
	'In case of sheet deletion
    Application.DisplayAlerts = False
    'It's already assured that this sheet exists - the custom handler works exclusively with this sheet
    With ActiveWorkbook.Sheets("MD_Config")
		'Parse string by string till cell is blank
        i = 1
        While .Cells(2 + i, 5).Value <> ""
			'Old dimension index
            OldIndex = Val(.Cells(2 + i, 5).Value)
            'New dimension index
            NewIndex = Val(.Cells(2 + i, 6).Value)
            'Check the bounds; if valid, start the sorting/renaming process
            If NewIndex >= 1 And NewIndex <= Count And OldIndex >= 1 And OldIndex <= Count Then PreSort OldIndex, NewIndex
            'Next string
            i = i + 1
        Wend
        'Finalize the sorting/renaming process
        For j = 1 To i - 1
			NewIndex = Val(.Cells(2 + j, 6).Value)
			If NewIndex >= 1 And NewIndex <= Count Then FinSort NewIndex
        Next
        'Move cursor
        .Cells(1, 2).Select
    End With
    'Restore default system value
    Application.DisplayAlerts = True
End Sub

'The handler procedure: treats the "LAUNCH" cell as a button that triggers the process of configuration update.
Private Sub sheetevent_SelectionChange(ByVal Target As Range)
	'We need the left top cell from the selected range 
    With Target.Cells(1, 1)
		'It's our button
        If .Value = "LAUNCH" And .Font.Size = 14 Then
			'Read the data from the configuration sheet: current (visible) dimension number and total count of dimensions
            VisibleDimension = Val(Cells(1, 2))
            NewCountOfDimensions = Val(Cells(2, 2))
            'Check the bounds and correct the values 
            If NewCountOfDimensions < 1 Then NewCountOfDimensions = 1
            If VisibleDimension > NewCountOfDimensions Then VisibleDimension = NewCountOfDimensions Else If VisibleDimension < 1 Then VisibleDimension = 1
            'Sorting/renaming process
            SortOut (NewCountOfDimensions)
            'Create the necessary number of dimensions
            CreateShadowPages (NewCountOfDimensions)
            'Erase the rest
            DestroyShadowPages (NewCountOfDimensions)
            'Store the total count of dimensions in the inner memory
            SetPageNow (NewCountOfDimensions)
            'Go to the dimension marked as visible
            SelectShadowPage (VisibleDimension)
            'Hide this configuration page
            ActiveWorkbook.Sheets("MD_Config").Visible = xlVeryHidden
        End If
    End With
End Sub
